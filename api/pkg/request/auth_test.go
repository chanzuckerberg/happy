package request

import (
	"context"
	"fmt"
	"testing"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
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
			app := fiber.New()
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			err := validateAuthHeader(ctx, tc.authHeader, tc.verifier)
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
			app := fiber.New()
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			err := validateAuthHeader(ctx, tc.authHeader, tc.verifier)
			r.Error(err)
		})
	}
}

func withSkipExpiryTokenCheck() providerVeriferOpt {
	return func(config *oidc.Config) {
		config.SkipExpiryCheck = true
	}
}
func TestGithubProvider(t *testing.T) {
	r := require.New(t)
	// The test tokens below will expire, but we still want to parse them and verify their tokens
	// were valid and the issuers claims match.
	verifier := MakeGithubVerifier("chanzuckerberg", withSkipExpiryTokenCheck())

	// NODE: this is not a real token anymore
	// it was at one point but has been expired
	// I am using it as a test so that I can make sure that we parse Github ID tokens properly
	// please don't yell at me :D
	tokens := []string{
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6ImVCWl9jbjNzWFlBZDBjaDRUSEJLSElnT3dPRSIsImtpZCI6Ijc4MTY3RjcyN0RFQzVEODAxREQxQzg3ODRDNzA0QTFDODgwRUMwRTEifQ.eyJqdGkiOiJlMDQ4ODBjZi1iZDFlLTQ5YmUtYjVkZi03MzNkOTEzZTljNWEiLCJzdWIiOiJyZXBvOmNoYW56dWNrZXJiZXJnL3NjcnVmZnk6cmVmOnJlZnMvaGVhZHMvaGVhdGhqL3ZpZXctaWQtdG9rZW4iLCJhdWQiOiJoYXBpLmhhcGkucHJvZC5zaS5jemkudGVjaG5vbG9neSIsInJlZiI6InJlZnMvaGVhZHMvaGVhdGhqL3ZpZXctaWQtdG9rZW4iLCJzaGEiOiIwNzVmMWVmMDdkMjA2NGU0MmFjMzM1N2FhNDIyOTA2OWMwYTRiOTUwIiwicmVwb3NpdG9yeSI6ImNoYW56dWNrZXJiZXJnL3NjcnVmZnkiLCJyZXBvc2l0b3J5X293bmVyIjoiY2hhbnp1Y2tlcmJlcmciLCJyZXBvc2l0b3J5X293bmVyX2lkIjoiMTk5MTcyOTkiLCJydW5faWQiOiI2MDE3NzkyMjAxIiwicnVuX251bWJlciI6IjMiLCJydW5fYXR0ZW1wdCI6IjEiLCJyZXBvc2l0b3J5X3Zpc2liaWxpdHkiOiJpbnRlcm5hbCIsInJlcG9zaXRvcnlfaWQiOiI2Njc2MDA5MTYiLCJhY3Rvcl9pZCI6Ijc2MDExOTEzIiwiYWN0b3IiOiJqYWtleWhlYXRoIiwid29ya2Zsb3ciOiJQcmludCBJRCB0b2tlbiIsImhlYWRfcmVmIjoiIiwiYmFzZV9yZWYiOiIiLCJldmVudF9uYW1lIjoicHVzaCIsInJlZl9wcm90ZWN0ZWQiOiJmYWxzZSIsInJlZl90eXBlIjoiYnJhbmNoIiwid29ya2Zsb3dfcmVmIjoiY2hhbnp1Y2tlcmJlcmcvc2NydWZmeS8uZ2l0aHViL3dvcmtmbG93cy9wcmludC1pZC10b2tlbi55bWxAcmVmcy9oZWFkcy9oZWF0aGovdmlldy1pZC10b2tlbiIsIndvcmtmbG93X3NoYSI6IjA3NWYxZWYwN2QyMDY0ZTQyYWMzMzU3YWE0MjI5MDY5YzBhNGI5NTAiLCJqb2Jfd29ya2Zsb3dfcmVmIjoiY2hhbnp1Y2tlcmJlcmcvc2NydWZmeS8uZ2l0aHViL3dvcmtmbG93cy9wcmludC1pZC10b2tlbi55bWxAcmVmcy9oZWFkcy9oZWF0aGovdmlldy1pZC10b2tlbiIsImpvYl93b3JrZmxvd19zaGEiOiIwNzVmMWVmMDdkMjA2NGU0MmFjMzM1N2FhNDIyOTA2OWMwYTRiOTUwIiwicnVubmVyX2Vudmlyb25tZW50Ijoic2VsZi1ob3N0ZWQiLCJlbnRlcnByaXNlIjoiY2hhbi16dWNrZXJiZXJnLWluaXRpYXRpdmUiLCJpc3MiOiJodHRwczovL3Rva2VuLmFjdGlvbnMuZ2l0aHVidXNlcmNvbnRlbnQuY29tIiwibmJmIjoxNjkzMzQ0NDAxLCJleHAiOjE2OTMzNDUzMDEsImlhdCI6MTY5MzM0NTAwMX0.IMZS-RwF4hL7ZXGIhjsAiP9v4s92mjYKv75feo0VZoN-5oQIp88cJRDEZvJmU-2pY5ZC5i_zxvPg-4DOlR936ZQss_8yMr9zayZsHG9bN-iP4H8_lvuaPdvUAJpCQNQj8AFYAMUOLP7lbHYlXFzfWgqzXCGmcZ-k6gMdHJBOSpkeIHoFsGg0VzmV9KeInwTKJE95ASrfxelXuOEcBjfeyF3tHAng7XXBx1Ls4v0xTZlbvhMqIcCZsCd-B5sQZF3yz2wP_84NAS_zf6E8QxoAs7VaEyG9uvJZwpEqN_1F3IBXJvjfbNmzo0Cun0BskxuSulOm6fBi2srzvqe8ICC3Xw",
	}

	for _, token := range tokens {
		idToken, err := verifier.Verify(context.Background(), token)
		r.NoError(err)
		r.NotNil(idToken)
	}
}
