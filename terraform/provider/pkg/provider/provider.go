package provider

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type PrivateKeyTFTokenProvider struct {
	signingKey              *rsa.PrivateKey
	issuer, audience, scope string
}

func MakeKMSKeyTFProvider(ctx context.Context, kmsKeyID, awsAssumeRoleARN, region, issuer, authzID, scope string) (client.TokenProvider, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load SDK config")
	}
	stsClient := sts.NewFromConfig(cfg)
	appCreds := stscreds.NewAssumeRoleProvider(stsClient, awsAssumeRoleARN)
	return &KMSKeyTFTokenProvider{
		client: kms.NewFromConfig(aws.Config{
			Credentials: appCreds,
			Region:      region,
		}),
		keyID:    kmsKeyID,
		issuer:   issuer,
		audience: fmt.Sprintf("https://czi-prod.okta.com/oauth2/%s/v1/token", authzID),
		scope:    scope,
	}, nil
}

type KMSKeyTFTokenProvider struct {
	client                         *kms.Client
	keyID, issuer, audience, scope string
}

func (k *KMSKeyTFTokenProvider) GetToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Issuer:    k.issuer,
		Subject:   k.issuer,
		Audience:  []string{k.audience},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	})
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

func MakePrivateKeyTFTokenProvider(rsaPrivateKey io.Reader, issuer, authzID, scope string) (client.TokenProvider, error) {
	b, err := io.ReadAll(rsaPrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "error reading private key")
	}
	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse RSA private key")
	}
	return &PrivateKeyTFTokenProvider{
		signingKey: signingKey,
		issuer:     issuer,
		audience:   fmt.Sprintf("https://czi-prod.okta.com/oauth2/%s/v1/token", authzID),
		scope:      scope,
	}, nil
}

type AccessTokenResponse struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}

func (t PrivateKeyTFTokenProvider) GetToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Issuer:    t.issuer,
		Subject:   t.issuer,
		Audience:  []string{t.audience},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	})
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

	accessTokenResp := AccessTokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&accessTokenResp)
	if err != nil {
		return "", errors.Wrap(err, "unable to decode access token response")
	}

	return accessTokenResp.AccessToken, nil
}

type APIClient struct {
	api client.HappyConfigAPI
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
				Required:      false,
				Description:   "The authentication credentials in the form of a PEM encoded private key to authenticate to the Happy API. Conflicts with api_kms_key_id.",
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
				DefaultFunc: schema.EnvDefaultFunc("HAPPY_API_OIDC_SCOPE", "scope"),
			},
			"api_kms_key_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Required:      false,
				Description:   "If set, the provider will use the KMS key ID to sign the JWT for the happy service user. The provider will need valid AWS credentials with access to the key. Conflicts with api_private_key.",
				DefaultFunc:   schema.EnvDefaultFunc("HAPPY_API_KMS_KEY_ID", "scope"),
				ConflictsWith: []string{"api_private_key"},
				RequiredWith:  []string{"api_assume_role_arn"},
			},
			"api_assume_role_arn": {
				Type:          schema.TypeString,
				Optional:      true,
				Required:      false,
				Description:   "The ARN of the role to assume when calling the KMS API to create a JWT signature.",
				DefaultFunc:   schema.EnvDefaultFunc("HAPPY_API_ASSUME_ROLE_ARN", "scope"),
				ConflictsWith: []string{"api_private_key"},
				RequiredWith:  []string{"api_kms_key_id"},
			},
			"api_kms_region": {
				Type:          schema.TypeString,
				Optional:      true,
				Required:      false,
				Description:   "The region the KMS key is located in. Defaults to us-west-2",
				DefaultFunc:   schema.EnvDefaultFunc("HAPPY_API_KMS_REGION", "us-west-2"),
				ConflictsWith: []string{"api_private_key"},
				RequiredWith:  []string{"api_kms_key_id"},
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"happy_resolved_app_configs": ResolvedAppConfigs(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiBaseURL := d.Get("api_base_url").(string)
	oidcIssuer := d.Get("api_oidc_issuer").(string)
	authzID := d.Get("api_oidc_authz_id").(string)
	scope := d.Get("api_oidc_scope").(string)

	var (
		tokenProvider client.TokenProvider
		err           error
	)
	if kmsKeyID, ok := d.GetOk("api_kms_key_id"); ok {
		if assumeRoleARN, ok := d.GetOk("api_assume_role_arn"); ok {
			tokenProvider, err = MakeKMSKeyTFProvider(ctx, kmsKeyID.(string), assumeRoleARN.(string), d.Get("api_kms_region").(string), oidcIssuer, authzID, scope)
			if err != nil {
				return nil, diag.FromErr(err)
			}
		}
	} else if apiPrivateKey, ok := d.GetOk("api_private_key"); ok {
		tokenProvider, err = MakePrivateKeyTFTokenProvider(strings.NewReader(apiPrivateKey.(string)), oidcIssuer, authzID, scope)
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	api := client.NewHappyClient("happy-provider", version.ProviderVersion, apiBaseURL, tokenProvider)
	return &APIClient{api: api}, nil
}
