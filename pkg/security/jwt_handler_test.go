package security

import (
	"encoding/base64"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

const (
	testAuthority = "course-watch"
	testAudience  = "course-watch-api"
	testTokenTtl  = time.Hour * 1
)

var testSigningKey = []byte("1234")

func TestBearerTokenClaims_IgnoresExpiration(t *testing.T) {
	// the expiration is validated by custom code
	issued := time.Now().Add(-2 * time.Hour) //issued 2 hours ago
	expires := issued.Add(1 * time.Hour)     //expired 1 hour ago
	claims := &bearerTokenClaims{}
	claims.NotBefore = jwt.NewNumericDate(issued)
	claims.IssuedAt = jwt.NewNumericDate(issued)
	claims.ExpiresAt = jwt.NewNumericDate(expires)
	err := claims.Valid()
	require.NoError(t, err)
}

func getReferenceJwtHandler() *JwtHandler {
	jh := NewJwtHandler()
	jh.Issuer = testAuthority
	jh.AudienceGenerated = []string{testAudience}
	jh.AudienceExpected = testAudience
	jh.TokenTtl = testTokenTtl
	jh.SigningKey = testSigningKey
	return jh
}

func getReferenceUser() *UserPrincipal {
	return &UserPrincipal{"1111111", []Role{Student}}
}

func decodeSegment(t *testing.T, seg string) []byte {
	t.Helper()
	result, err := base64.RawURLEncoding.DecodeString(seg)
	require.NoError(t, err)
	return result
}

func encodeSegment(bytes []byte) string {
	return base64.RawURLEncoding.EncodeToString(bytes)
}

func decodeClaims(t *testing.T, tokenString string) *bearerTokenClaims {
	t.Helper()

	parts := strings.Split(tokenString, ".")
	require.Equal(t, 3, len(parts))
	bytes := decodeSegment(t, parts[1])

	var btc bearerTokenClaims
	err := json.Unmarshal(bytes, &btc)
	require.NoError(t, err)
	return &btc
}

// returns a new token string with claims replaced by newClaims, without changing the signature
func replaceClaims(t *testing.T, tokenString string, newClaims *bearerTokenClaims) string {
	t.Helper()

	parts := strings.Split(tokenString, ".")
	require.Equal(t, 3, len(parts))

	bytes, err := json.Marshal(newClaims)
	require.NoError(t, err)
	parts[1] = encodeSegment(bytes)

	return strings.Join(parts, ".")
}

func generateModifiedReferenceTokenString(t *testing.T, claimsModifier func(*bearerTokenClaims)) string {
	t.Helper()

	jh := getReferenceJwtHandler()
	up := getReferenceUser()

	// This is identical to JwtHandler.Generate(), but with claimsModifier() applied before signing the token
	btc := jh.generateClaims(up)
	claimsModifier(btc)
	tokenString, err := jh.generateSignedString(btc)
	require.NoError(t, err)

	return tokenString
}

func applyOffset(ts *jwt.NumericDate, toAdd time.Duration) {
	result := ts.Time.Add(toAdd)
	ts.Time = result
}

// creates a claims modifier which offsets all timestamps in claims by toAdd
func offsetTimestamps(toAdd time.Duration) func(*bearerTokenClaims) {
	return func(btc *bearerTokenClaims) {
		applyOffset(btc.ExpiresAt, toAdd)
		applyOffset(btc.NotBefore, toAdd)
		applyOffset(btc.IssuedAt, toAdd)
	}
}

func TestTokenHandler_GenerateAndParseBack(t *testing.T) {
	jh := getReferenceJwtHandler()
	up := getReferenceUser()

	tokenString, err := jh.Generate(up)
	require.NoError(t, err)
	t.Log("Reference:", tokenString)

	payload, err := jh.Parse(tokenString)
	require.NoError(t, err)
	require.Equal(t, up, &payload.UserPrincipal)
	require.Equal(t, jh.Issuer, payload.Issuer)
	require.Equal(t, jh.AudienceGenerated, payload.Audience)
}

func TestTokenHandler_ParseDetectsInvalidSignature(t *testing.T) {
	jh := getReferenceJwtHandler()
	up := getReferenceUser()

	tokenString, err := jh.Generate(up)
	require.NoError(t, err)

	// Adding "admin" to roles and rebuilding the token, while using the old signature
	btc := decodeClaims(t, tokenString)
	require.Equal(t, []Role{Student}, btc.Roles)
	btc.Roles = append(btc.Roles, Admin)
	newTokenString := replaceClaims(t, tokenString, btc)

	// Invalid signature must be detected
	_, err = jh.Parse(newTokenString)
	require.Error(t, err)
	require.ErrorIs(t, err, jwt.ErrSignatureInvalid)

	t.Log("With bad signature:", newTokenString)

	// Invalid signature IS NOT detected when explicitly ignored
	payload, err := jh.ParseWithoutSignature(newTokenString)
	require.NoError(t, err)
	require.True(t, payload.IsAdmin())
}

// This is a classic attack: take a valid token, replace "alg" with "none" and drop the signature
func TestTokenHandler_ParseDetectsInvalidSignature_AlgNone(t *testing.T) {
	jh := getReferenceJwtHandler()
	up := getReferenceUser()

	btc := jh.generateClaims(up)

	// Adding "admin" to roles and rebuilding the token, while using the old signature
	require.Equal(t, []Role{Student}, btc.Roles)
	btc.Roles = append(btc.Roles, Admin)

	// jh.generateClaims() -> NewWithClaims() -> SignedString() is identical to JwtHandler.Generate()
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodNone, btc).SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	t.Log("With none method:", tokenString)

	// Method "none" must be detected
	_, err = jh.Parse(tokenString)
	require.Error(t, err)
	require.ErrorIs(t, err, jwt.ErrTokenSignatureInvalid)

	// Method "none" IS NOT detected when explicitly ignored
	payload, err := jh.ParseWithoutSignature(tokenString)
	require.NoError(t, err)
	require.True(t, payload.IsAdmin())
}

func TestTokenHandler_ParseDetectsInvalidSignature_DifferentKey(t *testing.T) {
	jh := getReferenceJwtHandler()
	up := getReferenceUser()

	jh.SigningKey = []byte("42-42-42")

	tokenString, err := jh.Generate(up)
	require.NoError(t, err)

	// the original key will be used for verification
	jh.SigningKey = testSigningKey

	// Invalid signature must be detected
	_, err = jh.Parse(tokenString)
	require.Error(t, err)
	require.ErrorIs(t, err, jwt.ErrSignatureInvalid)

	t.Log("Signature with a different key:", tokenString)

	// Invalid signature IS NOT detected when explicitly ignored
	_, err = jh.ParseWithoutSignature(tokenString)
	require.NoError(t, err)
}

func TestTokenHandler_ParseDetectsExpiration(t *testing.T) {
	// Token created 2 * TokenTtl ago, thus expired  1 TokenTtl ago
	tokenString := generateModifiedReferenceTokenString(t, offsetTimestamps(-2*testTokenTtl))

	jh := getReferenceJwtHandler()

	_, err := jh.Parse(tokenString)
	require.Error(t, err)
	require.ErrorIs(t, err, jwt.ErrTokenExpired)
	//t.Log(err)

	_, err = jh.ParseWithoutSignature(tokenString)
	require.Error(t, err)
	require.ErrorIs(t, err, jwt.ErrTokenExpired)
}

func TestTokenHandler_ParseDetectsNotBefore(t *testing.T) {
	// Token created 1 * TokenTtl in the future, thus not yet valid
	tokenString := generateModifiedReferenceTokenString(t, offsetTimestamps(1*testTokenTtl))

	jh := getReferenceJwtHandler()

	_, err := jh.Parse(tokenString)
	require.Error(t, err)
	require.ErrorIs(t, err, jwt.ErrTokenUsedBeforeIssued) // both must be present
	require.ErrorIs(t, err, jwt.ErrTokenNotValidYet)      // due to bit mask mechanism in the underlying error
	//t.Log(err)

	_, err = jh.ParseWithoutSignature(tokenString)
	require.Error(t, err)
	require.ErrorIs(t, err, jwt.ErrTokenUsedBeforeIssued) // both must be present
	require.ErrorIs(t, err, jwt.ErrTokenNotValidYet)      // due to bit mask mechanism in the underlying error
}

func TestTokenHandler_ParseDetectsMissingClaim_ExpiresAt(t *testing.T) {

	tokenString := generateModifiedReferenceTokenString(t, func(btc *bearerTokenClaims) {
		btc.ExpiresAt = nil
	})

	jh := getReferenceJwtHandler()

	token, err := jh.Parse(tokenString)
	require.True(t, token.ExpiresAt.IsZero())
	require.Error(t, err)
	require.ErrorIs(t, err, jwt.ErrTokenExpired)

	//t.Log(err)
	//t.Logf("%+v", token)

	token, err = jh.ParseWithoutSignature(tokenString)
	require.True(t, token.ExpiresAt.IsZero())
	require.Error(t, err)
	require.ErrorIs(t, err, jwt.ErrTokenExpired)
}

func TestTokenHandler_ParseDetectsMissingClaim_IssuedAt(t *testing.T) {

	tokenString := generateModifiedReferenceTokenString(t, func(btc *bearerTokenClaims) {
		btc.IssuedAt = nil
	})

	jh := getReferenceJwtHandler()

	token, err := jh.Parse(tokenString)
	require.True(t, token.IssuedAt.IsZero())
	require.Error(t, err)
	require.ErrorIs(t, err, jwt.ErrTokenUsedBeforeIssued)

	//t.Log(err)
	//t.Logf("%+v", token)

	token, err = jh.ParseWithoutSignature(tokenString)
	require.True(t, token.IssuedAt.IsZero())
	require.Error(t, err)
	require.ErrorIs(t, err, jwt.ErrTokenUsedBeforeIssued)
}

func TestTokenHandler_ParseDetectsMissingClaim_NotBefore(t *testing.T) {

	tokenString := generateModifiedReferenceTokenString(t, func(btc *bearerTokenClaims) {
		btc.NotBefore = nil
	})

	jh := getReferenceJwtHandler()

	token, err := jh.Parse(tokenString)
	require.True(t, token.NotBefore.IsZero())
	require.Error(t, err)
	require.ErrorIs(t, err, jwt.ErrTokenNotValidYet)

	//t.Log(err)
	//t.Logf("%+v", token)

	token, err = jh.ParseWithoutSignature(tokenString)
	require.True(t, token.NotBefore.IsZero())
	require.Error(t, err)
	require.ErrorIs(t, err, jwt.ErrTokenNotValidYet)
}

func TestTokenHandler_ParseDetectsMismatch_Issuer(t *testing.T) {

	generateTokenWithModifiedIssuer := func(t *testing.T, issuer string) string {
		t.Helper()

		jh := getReferenceJwtHandler()
		jh.Issuer = issuer
		up := getReferenceUser()

		tokenString, err := jh.Generate(up)
		require.NoError(t, err)
		return tokenString
	}

	cases := map[string]string{
		"different_issuer": generateTokenWithModifiedIssuer(t, "some-other-authority"),
		"empty_issuer":     generateTokenWithModifiedIssuer(t, ""),
	}

	for name, tokenString := range cases {
		t.Run(name, func(t *testing.T) {
			jh := getReferenceJwtHandler()

			_, err := jh.Parse(tokenString)
			require.Error(t, err)
			require.ErrorIs(t, err, jwt.ErrTokenInvalidIssuer)

			//t.Log(err)

			_, err = jh.ParseWithoutSignature(tokenString)
			require.Error(t, err)
			require.ErrorIs(t, err, jwt.ErrTokenInvalidIssuer)
		})
	}

}

func TestTokenHandler_ParseDetectsMismatch_Audience(t *testing.T) {

	generateTokenWithModifiedAudience := func(t *testing.T, audience []string) string {
		t.Helper()

		jh := getReferenceJwtHandler()
		jh.AudienceGenerated = audience
		up := getReferenceUser()

		tokenString, err := jh.Generate(up)
		require.NoError(t, err)
		return tokenString
	}

	cases := map[string]struct {
		tokenString string
		success     bool
	}{
		"multiple_single_match": {
			tokenString: generateTokenWithModifiedAudience(t, []string{testAudience, "some-other-api-2"}),
			success:     true,
		},
		"multiple_no_match": {
			tokenString: generateTokenWithModifiedAudience(t, []string{"some-other-api-1", "some-other-api-2"}),
			success:     false,
		},
		"different_audience": {
			tokenString: generateTokenWithModifiedAudience(t, []string{"some-other-api"}),
			success:     false,
		},
		"empty_audience": {
			tokenString: generateTokenWithModifiedAudience(t, []string{}),
			success:     false,
		},
		"nil_audience": {
			tokenString: generateTokenWithModifiedAudience(t, nil),
			success:     false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			jh := getReferenceJwtHandler()

			_, err := jh.Parse(tc.tokenString)

			if tc.success {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorIs(t, err, jwt.ErrTokenInvalidAudience)
			}

			_, err = jh.ParseWithoutSignature(tc.tokenString)
			if tc.success {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorIs(t, err, jwt.ErrTokenInvalidAudience)
			}
		})
	}
}
