package request

import (
	"context"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
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

// TODO: add more claims here based on what we expect
// we can use these to grant access to specific stacks and resources
type claims struct {
	Email   string `json:"email"`
	Subject string `json:"sub"`
}

type OIDCVerifier interface {
	Verify(context.Context, string) (*oidc.IDToken, error)
}

func MakeOIDCVerifier(ctx context.Context, issuerURL, clientID string) (OIDCVerifier, error) {
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create OIDC provider from %s", issuerURL)
	}

	return provider.Verifier(&oidc.Config{ClientID: clientID}), nil
}

func validateAuthHeader(ctx context.Context, authHeader string, verifier OIDCVerifier) (string, string, error) {
	rawIDToken := stripBearerPrefixFromTokenString(authHeader)
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return "", "", errors.Wrap(err, "unable to verify ID token")
	}

	claims := &claims{}
	err = idToken.Claims(claims)
	if err != nil {
		return "", "", errors.Wrap(err, "unable to parse claims from ID token")
	}
	if len(claims.Email) == 0 {
		return "", "", errors.New("missing email claim")
	}
	if len(claims.Subject) == 0 {
		return "", "", errors.New("missing subject claim")
	}

	return claims.Email, claims.Subject, nil
}

func MakeAuth(verifier OIDCVerifier) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.GetRespHeader("Authorization")
		if len(authHeader) <= 0 {
			return c.SendStatus(fiber.StatusForbidden)
		}

		email, subject, err := validateAuthHeader(c.Context(), authHeader, verifier)
		if err != nil {
			return c.SendStatus(fiber.StatusForbidden)
		}

		c.Locals("email", email)
		c.Locals("subject", subject)
		return c.Next()
	}
}
