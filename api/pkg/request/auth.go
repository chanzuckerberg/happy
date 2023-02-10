package request

import (
	"context"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// From https://github.com/dgrijalva/jwt-go/blob/master/request/oauth2.go
// Strips 'Bearer ' prefix from bearer token string
func stripBearerPrefixFromTokenString(token string) string {
	// Should be a bearer token
	if len(token) > 6 && strings.ToUpper(token[0:7]) == "BEARER " {
		return token[7:]
	}
	return token
}

type OIDCVerifier interface {
	Verify(ctx context.Context, idToken string) (*oidc.IDToken, error)
}

type ClaimsVerifier interface {
	MatchClaims(ctx context.Context, idToken *oidc.IDToken) error
}

type OIDCProvider struct {
	oidcVerifier   OIDCVerifier
	claimsVerifier ClaimsVerifier
}

type GithubClaimsVerifier struct {
	owner  string
	issuer string
}

type GithubClaims struct {
	Subject         string `json:"sub"`
	Issuer          string `json:"iss"`
	RepositoryOwner string `json:"repository_owner"`
	Repository      string `json:"repository"`
	Action          string `json:"actor"`
	HeadRef         string `json:"head_ref"`
	WorkflowSHA     string `json:"workflow_sha"`
}

type NilClaimsVerifier struct {
}

func (d *NilClaimsVerifier) MatchClaims(ctx context.Context, idToken *oidc.IDToken) error {
	return nil
}

var DefaultClaimsVerifier = &NilClaimsVerifier{}

func MakeGithubClaimsVerifier(owner string) *GithubClaimsVerifier {
	return &GithubClaimsVerifier{
		owner:  owner,
		issuer: "https://token.actions.githubusercontent.com",
	}
}

func (g *GithubClaimsVerifier) MatchClaims(ctx context.Context, idToken *oidc.IDToken) error {
	claims := GithubClaims{}
	err := idToken.Claims(&claims)
	if err != nil {
		return errors.Wrap(err, "github id token didn't have expected claims")
	}

	if claims.RepositoryOwner != g.owner {
		return errors.Errorf("github id token didn't have the expected github owner, expected %s got %s", g.owner, claims.RepositoryOwner)
	}

	if claims.Issuer != g.issuer {
		return errors.Errorf("github id token didn't have the expected issuer, expected %s got %s", g.issuer, claims.Issuer)
	}

	return nil
}

func (o *OIDCProvider) Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
	idToken, err := o.oidcVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't verify the ID token")
	}

	err = o.claimsVerifier.MatchClaims(ctx, idToken)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't verify the claims of the ID token")
	}

	return idToken, nil
}

type MultiOIDCProvider struct {
	verifiers []OIDCVerifier
}

func (m *MultiOIDCProvider) Verify(ctx context.Context, idToken string) (*oidc.IDToken, error) {
	var errs error
	for _, verifier := range m.verifiers {
		idToken, err := verifier.Verify(ctx, idToken)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
		return idToken, nil
	}
	return nil, errors.Wrap(errs, "unable to verify the ID token with any of the configured OIDC providers")
}

func MakeMultiOIDCVerifier(verifiers ...OIDCVerifier) OIDCVerifier {
	return &MultiOIDCProvider{
		verifiers: verifiers,
	}
}

func MakeOIDCProvider(ctx context.Context, issuerURL, clientID string, claimsVerifier ClaimsVerifier) (*OIDCProvider, error) {
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create OIDC provider from %s", issuerURL)
	}

	return &OIDCProvider{
		oidcVerifier:   provider.Verifier(&oidc.Config{ClientID: clientID}),
		claimsVerifier: claimsVerifier,
	}, nil
}

func validateAuthHeader(ctx context.Context, authHeader string, verifier OIDCVerifier) error {
	rawIDToken := stripBearerPrefixFromTokenString(authHeader)
	_, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return errors.Wrap(err, "unable to verify ID token")
	}
	// TODO: once we have some common patterns of access, extra these properties
	// from the ID token here and attach them to the request using
	// fiber.Ctx.Locals(key, value)
	return nil
}

func MakeAuth(verifier OIDCVerifier) fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()
		authHeader := headers["Authorization"]
		if len(authHeader) <= 0 {
			return c.SendStatus(fiber.StatusForbidden)
		}
		err := validateAuthHeader(c.Context(), authHeader, verifier)
		if err != nil {
			return c.SendStatus(fiber.StatusForbidden)
		}
		return c.Next()
	}
}
