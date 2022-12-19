package provider

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/chanzuckerberg/happy/shared/client"
	"github.com/chanzuckerberg/happy/terraform/provider/pkg/version"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

type TFTokenProvider struct {
	signingKey              *rsa.PrivateKey
	issuer, audience, scope string
}

func MakeTFTokenProvider(rsaPrivateKey io.Reader, issuer, authzID, scope string) (*TFTokenProvider, error) {
	b, err := io.ReadAll(rsaPrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "error reading private key")
	}
	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse RSA private key")
	}
	return &TFTokenProvider{
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

func (t TFTokenProvider) GetToken() (string, error) {
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

	values := url.Values{}
	values.Add("grant_type", "client_credentials")
	values.Add("scope", t.scope)
	values.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	values.Add("client_assertion", signedToken)

	params := values.Encode()
	req, err := http.NewRequest(http.MethodPost, t.audience, strings.NewReader(params))
	if err != nil {
		return "", errors.Wrapf(err, "error talking %s", t.audience)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "error talking to %s", t.audience)
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "The authentication credentials in the form of a PEM encoded private key to authenticate to the Happy API.",
				DefaultFunc: schema.EnvDefaultFunc("HAPPY_API_TOKEN", nil),
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
	apiPrivateKey := d.Get("api_private_key").(string)
	oidcIssuer := d.Get("api_oidc_issuer").(string)
	authzID := d.Get("api_oidc_authz_id").(string)
	scope := d.Get("api_oidc_scope").(string)

	provider, err := MakeTFTokenProvider(strings.NewReader(apiPrivateKey), oidcIssuer, authzID, scope)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	api := client.NewHappyClient("happy-provider", version.ProviderVersion, apiBaseURL, provider)
	return &APIClient{api: api}, nil
}
