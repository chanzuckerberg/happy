package config

const (
	DEFAULT_HAPPY_API_BASE_URL        = "https://hapi.hapi.prod.si.czi.technology"
	DEFAULT_HAPPY_API_OIDC_CLIENT_ID  = "0oa8anwuhpAX1rfvb5d7"
	DEFAULT_HAPPY_API_OIDC_ISSUER_URL = "https://czi-prod.okta.com"
)

type HappyApiConfig struct {
	BaseUrl       string `yaml:"base_url"`
	OidcClientID  string `yaml:"oidc_client_id"`
	OidcIssuerUrl string `yaml:"oidc_issuer_url"`
}

func DefaultHappyApiConfig() HappyApiConfig {
	return HappyApiConfig{
		BaseUrl:       DEFAULT_HAPPY_API_BASE_URL,
		OidcClientID:  DEFAULT_HAPPY_API_OIDC_CLIENT_ID,
		OidcIssuerUrl: DEFAULT_HAPPY_API_OIDC_ISSUER_URL,
	}
}
