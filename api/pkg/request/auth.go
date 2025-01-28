package request

import (
	"context"
	"strings"

	"github.com/chanzuckerberg/happy/api/pkg/ent/ogent"
	"github.com/chanzuckerberg/happy/api/pkg/response"
	"github.com/chanzuckerberg/happy/api/pkg/setup"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-multierror"
	"github.com/ogen-go/ogen/middleware"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

type GithubVerifier struct {
	opts []providerVeriferOpt
	OIDCProvider
}

type providerVeriferOpt func(*oidc.Config)

func (g *GithubVerifier) Verify(ctx context.Context, idToken string) (*oidc.IDToken, error) {
	provider, err := oidc.NewProvider(ctx, "https://token.actions.githubusercontent.com")
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create github oidc provider")
	}

	config := &oidc.Config{
		SkipClientIDCheck: true,
	}
	for _, opt := range g.opts {
		opt(config)
	}
	verifier := provider.Verifier(config)
	return verifier.Verify(ctx, idToken)
}

func MakeGithubVerifier(githubOwner string, opts ...providerVeriferOpt) *GithubVerifier {
	return &GithubVerifier{
		opts: opts,
		OIDCProvider: OIDCProvider{
			oidcVerifier:   &GithubVerifier{},
			claimsVerifier: MakeGithubClaimsVerifier(githubOwner),
		},
	}
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
	errs := &multierror.Error{}
	for _, verifier := range m.verifiers {
		idToken, err := verifier.Verify(ctx, idToken)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
		return idToken, nil
	}
	return nil, errors.Wrap(errs.ErrorOrNil(), "unable to verify the ID token with any of the configured OIDC providers")
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

	config := &oidc.Config{ClientID: clientID}
	return &OIDCProvider{
		oidcVerifier:   provider.Verifier(config),
		claimsVerifier: claimsVerifier,
	}, nil
}

type OIDCAuthKey struct{}

type OIDCAuthValues struct {
	Subject string
	Email   string
	Actor   string
}

func ValidateAuthHeader(ctx context.Context, authHeader string, verifier OIDCVerifier) (*OIDCAuthValues, error) {
	rawIDToken := stripBearerPrefixFromTokenString(authHeader)
	token, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, errors.Wrap(err, "unable to verify ID token")
	}

	var claims struct {
		Email string `json:"email"`
		Actor string `json:"actor"`
	}
	err = token.Claims(&claims)
	if err != nil {
		return nil, err
	}
	if claims.Email == "" && claims.Actor == "" {
		// TODO: can't throw an error here because it breaks TFE runs, log the issue for now
		// return errors.New("ID token didn't have email or actor claims")
		logrus.Warn("ID token didn't have email or actor claims")
	}

	return &OIDCAuthValues{
		Subject: token.Subject,
		Email:   claims.Email,
		Actor:   claims.Actor,
	}, nil
}

func MakeVerifierFromConfig(ctx context.Context, cfg *setup.Configuration) OIDCVerifier {
	verifiers := []OIDCVerifier{
		MakeGithubVerifier("chanzuckerberg"),
	}
	for _, provider := range cfg.Auth.Providers {
		verifier, err := MakeOIDCProvider(ctx, provider.IssuerURL, provider.ClientID, DefaultClaimsVerifier)
		if err != nil {
			logrus.Errorf("failed to create OIDC verifier with error: %s", err.Error())
			continue
		}
		verifiers = append(verifiers, verifier)
	}

	return MakeMultiOIDCVerifier(verifiers...)
}

func MakeFiberAuthMiddleware(verifier OIDCVerifier) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.GetReqHeaders()[fiber.HeaderAuthorization]
		if len(authHeader) <= 0 {
			return response.AuthErrorResponse(c, "missing auth header")
		}

		oidcValues, err := ValidateAuthHeader(c.Context(), authHeader[0], verifier)
		if err != nil {
			return response.AuthErrorResponse(c, err.Error())
		}

		c.Locals(OIDCAuthKey{}, oidcValues)
		return c.Next()
	}
}

func MakeOgentAuthMiddleware(verifier OIDCVerifier) ogent.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		authHeader := req.Raw.Header.Get("Authorization")
		if len(authHeader) <= 0 {
			return middleware.Response{}, response.NewForbiddenError("missing auth header")
		}

		oidcValues, err := ValidateAuthHeader(req.Context, authHeader, verifier)
		if err != nil {
			return middleware.Response{}, response.NewForbiddenError("you are not allowed to access this resource")
		}

		req.Context = context.WithValue(req.Context, OIDCAuthKey{}, oidcValues)
		user := sentry.User{}
		if len(oidcValues.Email) > 0 {
			user.Email = oidcValues.Email
		}
		if len(oidcValues.Actor) > 0 {
			user.Username = oidcValues.Actor
		}
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetUser(user)
		})

		return next(req)
	}
}
