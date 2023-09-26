package provider

import (
	"context"
	"os"
	"testing"

	"github.com/chanzuckerberg/happy/shared/client"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type APIMock struct {
	client.HappyConfigAPI
	mock.Mock
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func getTestProviders() (map[string]*schema.Provider, *APIMock) {
	happyProvider := Provider()
	apiMock := &APIMock{}
	happyProvider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		_, err := validateConfiguration(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		client := &APIClient{api: apiMock}
		return client, nil
	}
	providers := map[string]*schema.Provider{
		"happy": happyProvider,
	}
	return providers, apiMock
}

func testPreCheck(t *testing.T) {
	if err := os.Getenv("HAPPY_API_BASE_URL"); err == "" {
		t.Fatal("HAPPY_API_BASE_URL must be set for acceptance tests")
	}
}

func TestGetClaimsValid(t *testing.T) {
	claimsValues := ClaimsValues{
		issuer:        "issuer-value",
		audience:      "audience-value",
		scope:         "scope-value",
		assumeRoleARN: "arn:aws:iam::1234567890:role/fake-role",
	}
	claims := getClaims(claimsValues)

	r := require.New(t)
	r.Equal(claims.Audience, jwt.ClaimStrings(jwt.ClaimStrings{claimsValues.audience}))
	r.Equal(claims.Issuer, claimsValues.issuer)
	r.Equal(claims.Actor, claimsValues.assumeRoleARN)
	r.Equal(claims.Valid(), nil)
}

func TestGetClaimsInalid(t *testing.T) {
	claimsValues := ClaimsValues{
		issuer:        "issuer-value",
		audience:      "audience-value",
		scope:         "scope-value",
		assumeRoleARN: "",
	}
	claims := getClaims(claimsValues)

	r := require.New(t)
	r.Equal(claims.Audience, jwt.ClaimStrings(jwt.ClaimStrings{claimsValues.audience}))
	r.Equal(claims.Issuer, claimsValues.issuer)
	r.Equal(claims.Actor, claimsValues.assumeRoleARN)
	r.NotEqual(claims.Valid(), nil)
}
