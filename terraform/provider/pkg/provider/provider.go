package provider

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/chanzuckerberg/happy/shared/client"
	"github.com/chanzuckerberg/happy/terraform/provider/pkg/version"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

type ClaimsValues struct {
	issuer, audience, scope, assumeRoleARN string
}

type PrivateKeyTFTokenProvider struct {
	signingKey *rsa.PrivateKey
	ClaimsValues
}

func MakeKMSKeyTFProvider(ctx context.Context, provConfig *Config, appCreds *stscreds.AssumeRoleProvider) (client.TokenProvider, error) {
	return &KMSKeyTFTokenProvider{
		client: kms.NewFromConfig(aws.Config{
			Credentials: appCreds,
			Region:      provConfig.Region,
		}),
		keyID: *provConfig.KMSKeyID,
		ClaimsValues: ClaimsValues{
			issuer:        provConfig.OIDCIssuer,
			audience:      fmt.Sprintf("https://czi-prod.okta.com/oauth2/%s/v1/token", provConfig.OIDCAuthzID),
			scope:         provConfig.OIDCScope,
			assumeRoleARN: provConfig.AssumeRoleARN,
		},
	}, nil
}

type KMSKeyTFTokenProvider struct {
	client *kms.Client
	keyID  string
	ClaimsValues
}

type ExtendedClaims struct {
	jwt.RegisteredClaims
	Actor string `json:"actor,omitempty"`
}

func (c ExtendedClaims) Valid() error {
	if err := c.RegisteredClaims.Valid(); err != nil {
		return err
	}
	if c.Actor == "" {
		return errors.New("actor is required")
	}

	return nil
}

func getClaims(values ClaimsValues) ExtendedClaims {
	return ExtendedClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    values.issuer,
			Subject:   values.issuer,
			Audience:  []string{values.audience},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Actor: values.assumeRoleARN,
	}
}

func (k *KMSKeyTFTokenProvider) GetToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, getClaims(k.ClaimsValues))
	signingStr, err := token.SigningString()
	if err != nil {
		return "", errors.Wrap(err, "unable to make signing string")
	}

	signResponse, err := k.client.Sign(context.Background(), &kms.SignInput{
		Message:          []byte(signingStr),
		KeyId:            aws.String(k.keyID),
		SigningAlgorithm: types.SigningAlgorithmSpecRsassaPkcs1V15Sha256,
		MessageType:      types.MessageTypeRaw,
	})
	if err != nil {
		return "", errors.Wrapf(err, "unable to sign JWT with KMS key %s", k.keyID)
	}

	return requestAccessToken(k.scope, k.audience, fmt.Sprintf("%s.%s", signingStr, base64.RawStdEncoding.EncodeToString(signResponse.Signature)))
}

func MakePrivateKeyTFTokenProvider(provConfig *Config) (client.TokenProvider, error) {
	b, err := io.ReadAll(provConfig.PrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "error reading private key")
	}
	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse RSA private key")
	}
	return &PrivateKeyTFTokenProvider{
		signingKey: signingKey,
		ClaimsValues: ClaimsValues{
			issuer:        provConfig.OIDCIssuer,
			audience:      fmt.Sprintf("https://czi-prod.okta.com/oauth2/%s/v1/token", provConfig.OIDCAuthzID),
			scope:         provConfig.OIDCScope,
			assumeRoleARN: provConfig.AssumeRoleARN,
		},
	}, nil
}

type AccessTokenResponse struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}

func (t PrivateKeyTFTokenProvider) GetToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, getClaims(t.ClaimsValues))
	signedToken, err := token.SignedString(t.signingKey)
	if err != nil {
		return "", errors.Wrapf(err, "error signing the JWT for %s %s", t.issuer, t.audience)
	}

	return requestAccessToken(t.scope, t.audience, signedToken)
}

func requestAccessToken(scope, audience, signedToken string) (string, error) {
	values := url.Values{}
	values.Add("grant_type", "client_credentials")
	values.Add("scope", scope)
	values.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	values.Add("client_assertion", signedToken)

	params := values.Encode()
	req, err := http.NewRequest(http.MethodPost, audience, strings.NewReader(params))
	if err != nil {
		return "", errors.Wrapf(err, "error talking %s", audience)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "error talking to %s", audience)
	}
	if resp.StatusCode >= 300 {
		respOut, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return "", errors.Wrapf(err, "got status code %s %+v", resp.Status, resp.StatusCode)
		}
		return "", errors.Wrapf(err, "got status code %s %+v %s", resp.Status, resp.StatusCode, string(respOut))
	}

	accessTokenResp := AccessTokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&accessTokenResp)
	if err != nil {
		return "", errors.Wrap(err, "unable to decode access token response")
	}

	return accessTokenResp.AccessToken, nil
}

type AWSCredentialsProviderTF struct {
	appCreds *stscreds.AssumeRoleProvider
}

func (t AWSCredentialsProviderTF) GetCredentials(ctx context.Context) (aws.Credentials, error) {
	return t.appCreds.Retrieve(ctx)
}

func getAwsCredsProvider(appCreds *stscreds.AssumeRoleProvider) client.AWSCredentialsProvider {
	return AWSCredentialsProviderTF{
		appCreds: appCreds,
	}
}

type HappyAPIClient interface {
	client.HappyConfigAPI
	client.HappyStackAPI
}
type APIClient struct {
	api HappyAPIClient
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_base_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The base URL for where the Happy API is located.",
				DefaultFunc: schema.EnvDefaultFunc("HAPPY_API_BASE_URL", nil),
			},
			"api_private_key": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The authentication credentials in the form of a PEM encoded private key to authenticate to the Happy API.",
				DefaultFunc:   schema.EnvDefaultFunc("HAPPY_API_PRIVATE_KEY", nil),
				ConflictsWith: []string{"api_kms_key_id"},
			},
			"api_oidc_issuer": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The OIDC issuer value that corresponds to the client ID of the Okta application associated with the private key.",
				DefaultFunc: schema.EnvDefaultFunc("HAPPY_API_OIDC_ISSUER", nil),
			},
			"api_oidc_authz_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Okta authorization server ID that authenticates the requests to the Happy API.",
				DefaultFunc: schema.EnvDefaultFunc("HAPPY_API_OIDC_AUTHZ_ID", nil),
			},
			"api_oidc_scope": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The required scope for the service account to authenticate properly.",
				DefaultFunc: schema.EnvDefaultFunc("HAPPY_API_OIDC_SCOPE", nil),
			},
			"api_kms_key_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "If set, the provider will use the KMS key ID to sign the JWT for the happy service user. The provider will need valid AWS credentials via an assume role (set in api_assume_role_arn) with access to the key.",
				DefaultFunc:   schema.EnvDefaultFunc("HAPPY_API_KMS_KEY_ID", nil),
				ConflictsWith: []string{"api_private_key"},
			},
			"api_assume_role_arn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ARN of the role to assume when calling the KMS API to create a JWT signature.",
				DefaultFunc: schema.EnvDefaultFunc("HAPPY_API_ASSUME_ROLE_ARN", nil),
			},
			"api_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The AWS region to use when calling KMS.",
				DefaultFunc: schema.EnvDefaultFunc("HAPPY_API_REGION", "us-west-2"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"happy_resolved_app_configs": ResolvedAppConfigs(),
			"happy_stacklist":            Stacklist(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config, err := validateConfiguration(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	appCreds, err := getAppCreds(ctx, config)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	tokenProvider, err := getTokenProvider(ctx, config, appCreds)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	awsCredsProvider := getAwsCredsProvider(appCreds)

	return &APIClient{
		api: client.NewHappyClient("happy-provider", version.ProviderVersion, config.BaseURL, tokenProvider, awsCredsProvider),
	}, nil
}

func getAppCreds(ctx context.Context, provConfig *Config) (*stscreds.AssumeRoleProvider, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(provConfig.Region))
	if err != nil {
		return nil, errors.Wrap(err, "unable to load SDK config")
	}
	stsClient := sts.NewFromConfig(cfg)
	return stscreds.NewAssumeRoleProvider(stsClient, provConfig.AssumeRoleARN), nil
}

type Config struct {
	BaseURL       string
	OIDCIssuer    string
	OIDCAuthzID   string
	OIDCScope     string
	AssumeRoleARN string
	Region        string
	KMSKeyID      *string
	PrivateKey    io.Reader
}

func getTokenProvider(ctx context.Context, config *Config, appCreds *stscreds.AssumeRoleProvider) (client.TokenProvider, error) {
	if config.KMSKeyID != nil {
		return MakeKMSKeyTFProvider(ctx, config, appCreds)
	} else if config.PrivateKey != nil {
		return MakePrivateKeyTFTokenProvider(config)
	}
	return nil, errors.New("invalid configuration to get token provider")
}

func validateConfiguration(d *schema.ResourceData) (*Config, error) {
	config := &Config{
		BaseURL:       d.Get("api_base_url").(string),
		OIDCIssuer:    d.Get("api_oidc_issuer").(string),
		OIDCAuthzID:   d.Get("api_oidc_authz_id").(string),
		OIDCScope:     d.Get("api_oidc_scope").(string),
		AssumeRoleARN: d.Get("api_assume_role_arn").(string),
		Region:        d.Get("api_region").(string),
	}
	kmsKeyID, ok := d.GetOk("api_kms_key_id")
	if ok {
		config.KMSKeyID = aws.String(kmsKeyID.(string))
	}

	apiPrivateKey, ok := d.GetOk("api_private_key")
	if ok {
		config.PrivateKey = strings.NewReader(apiPrivateKey.(string))
	}

	if config.KMSKeyID == nil && config.PrivateKey == nil {
		return nil, errors.New("bad happy provider configuration, need KMS key ID (HAPPY_API_KMS_KEY_ID) or a happy API private key (HAPPY_API_PRIVATE_KEY)")
	}

	return config, nil
}
