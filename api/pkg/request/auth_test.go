package request

import (
	"fmt"
	"testing"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestStripBearer(t *testing.T) {
	r := require.New(t)
	type testCase struct {
		in, out string
	}
	testCases := []testCase{
		{in: "bearer blah", out: "blah"},
		{in: "Bearer blah", out: "blah"},
		{in: "BEARER blah", out: "blah"},
		{in: "BEArER blah", out: "blah"},
		{in: "bearerblah", out: "bearerblah"},
		{in: "Bearerszxcvasdf blah", out: "Bearerszxcvasdf blah"},
		{in: "blah", out: "blah"},
		{in: "", out: ""},
	}
	for _, testcase := range testCases {
		tc := testcase
		t.Run(tc.in, func(t *testing.T) {
			t.Parallel()
			out := stripBearerPrefixFromTokenString(tc.in)
			r.Equal(tc.out, out)
		})
	}
}

func newDummyJWT(r *require.Assertions, subject, email string) string {
	mapClaims := jwt.MapClaims{
		"sub":   subject,
		"email": email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	ss, err := token.SignedString([]byte{})
	r.NoError(err)
	return ss
}

func TestValidateAuthHeaderNoErrors(t *testing.T) {
	r := require.New(t)
	type testCase struct {
		authHeader      string
		verifier        OIDCVerifier
		expectedSubject string
	}
	dummyVerifier := oidc.NewVerifier("blah", nil, &oidc.Config{
		SkipClientIDCheck:          true,
		SkipIssuerCheck:            true,
		SkipExpiryCheck:            true,
		InsecureSkipSignatureCheck: true,
	})
	testCases := []testCase{
		{
			authHeader:      fmt.Sprintf("Bearer %s", newDummyJWT(r, "subject", "exp@example.com")),
			verifier:        dummyVerifier,
			expectedSubject: "subject",
		},
		{
			authHeader:      newDummyJWT(r, "subject", "exp@example.com"), // Bearer is optional
			verifier:        dummyVerifier,
			expectedSubject: "subject",
		},
	}
	for _, testcase := range testCases {
		tc := testcase
		t.Run(tc.authHeader, func(t *testing.T) {
			t.Parallel()
			err := validateAuthHeader(&fiber.Ctx{}, tc.authHeader, tc.verifier)
			r.NoError(err)
		})
	}
}

func TestValidateAuthHeaderErrors(t *testing.T) {
	r := require.New(t)
	type testCase struct {
		authHeader string
		verifier   OIDCVerifier
	}
	dummyVerifier := oidc.NewVerifier("blah", nil, &oidc.Config{
		SkipClientIDCheck:          true,
		SkipIssuerCheck:            true,
		SkipExpiryCheck:            true,
		InsecureSkipSignatureCheck: true,
	})
	testCases := []testCase{
		{
			authHeader: fmt.Sprintf("Bearer %s", "blah"), // malformed JWT
			verifier:   dummyVerifier,
		},
	}
	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.authHeader, func(t *testing.T) {
			t.Parallel()
			err := validateAuthHeader(&fiber.Ctx{}, tc.authHeader, tc.verifier)
			r.Error(err)
		})
	}
}
