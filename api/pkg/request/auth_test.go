package request

import (
	"context"
	"fmt"
	"testing"

	"github.com/coreos/go-oidc/v3/oidc"
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

func newDummyJWTNoClaims(r *require.Assertions) string {
	token := jwt.New(jwt.SigningMethodHS256)
	ss, err := token.SignedString([]byte{})
	r.NoError(err)
	return ss
}

func TestValidateAuthHeaderNoErrors(t *testing.T) {
	r := require.New(t)
	type testCase struct {
		authHeader                     string
		verifier                       OIDCVerifier
		expectedEmail, expectedSubject string
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
			expectedEmail:   "exp@example.com",
			expectedSubject: "subject",
		},
		{
			authHeader:      newDummyJWT(r, "subject", "exp@example.com"), // Bearer is optional
			verifier:        dummyVerifier,
			expectedEmail:   "exp@example.com",
			expectedSubject: "subject",
		},
	}
	for _, testcase := range testCases {
		tc := testcase
		t.Run(tc.authHeader, func(t *testing.T) {
			t.Parallel()
			email, subject, err := validateAuthHeader(context.Background(), tc.authHeader, tc.verifier)
			r.NoError(err)
			r.Equal(tc.expectedEmail, email)
			r.Equal(tc.expectedSubject, subject)
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
		{
			authHeader: newDummyJWTNoClaims(r), //missing claims
			verifier:   dummyVerifier,
		},
	}
	for _, test := range testCases {
		t.Run(test.authHeader, func(t *testing.T) {
			t.Parallel()
			email, subject, err := validateAuthHeader(context.Background(), test.authHeader, test.verifier)
			r.Error(err)
			r.Equal("", email)
			r.Equal("", subject)
		})
	}
}
