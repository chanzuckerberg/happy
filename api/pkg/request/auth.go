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

type MultiOIDCVerifier struct {
	verifiers []OIDCVerifier
}

func (m *MultiOIDCVerifier) Verify(ctx context.Context, idToken string) (*oidc.IDToken, error) {
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
	return &MultiOIDCVerifier{
		verifiers: verifiers,
	}
}

func MakeOIDCVerifier(ctx context.Context, issuerURL, clientID string) (OIDCVerifier, error) {
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create OIDC provider from %s", issuerURL)
	}

	return provider.Verifier(&oidc.Config{ClientID: clientID}), nil
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
		authHeader := c.GetReqHeaders()["Authorization"]
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
