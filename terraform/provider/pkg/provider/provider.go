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
	signingKey                           *rsa.PrivateKey
	issuer, authorizationServerID, scope string
}

func MakeTFTokenProvider(rsaPrivateKey io.Reader, issuer, authorizationServerID, scope string) (*TFTokenProvider, error) {
	b, err := io.ReadAll(rsaPrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "error reading private key")
	}
	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse RSA private key")
	}
	return &TFTokenProvider{
		signingKey:            signingKey,
		issuer:                issuer,
		authorizationServerID: authorizationServerID,
		scope:                 scope,
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
		Audience:  []string{fmt.Sprintf("https://czi-prod.okta.com/oauth2/%s/v1/token", t.authorizationServerID)},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	})
	signedToken, err := token.SignedString(t.signingKey)
	if err != nil {
		return "", errors.Wrapf(err, "error signing the JWT for %s %s", t.issuer, t.authorizationServerID)
	}

	values := url.Values{}
	values.Add("grant_type", "client_credentials")
	values.Add("scope", t.scope)
	values.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	values.Add("client_assertion", signedToken)

	url := fmt.Sprintf("https://czi-prod.okta.com/oauth2/%s/v1/token", t.authorizationServerID)
	params := values.Encode()
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(params))
	if err != nil {
		return "", errors.Wrapf(err, "error talking %s", url)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "error talking %s", url)
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
				DefaultFunc: schema.EnvDefaultFunc("HAPPY_API_BASE_URL", nil),
			},
			// TODO: add inputs needed for oidc
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"happy_resolved_app_configs": ResolvedAppConfigs(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiBaseUrl := d.Get("api_base_url").(string)
	api := client.NewHappyClient("happy-provider", version.ProviderVersion, apiBaseUrl, TFTokenProvider{})
	return &APIClient{api: api}, nil
}
