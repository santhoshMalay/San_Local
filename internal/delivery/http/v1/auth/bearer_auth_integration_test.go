package auth_test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zhuravlev-pe/course-watch/internal/delivery/http/v1/auth"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var validKey = []byte("1234")
var invalidKey = []byte("42-42-42")

func getParametrizedTokenHandler(key []byte) *security.JwtHandler {
	jwt := security.NewJwtHandler()

	const aud = "course-watch-api"

	jwt.Issuer = "course-watch"
	jwt.AudienceExpected = aud
	jwt.AudienceGenerated = []string{aud}
	jwt.SigningKey = key
	jwt.TokenTtl = time.Hour * 1

	return jwt
}

func TestBearerAuthenticator_Integration_Authorize(t *testing.T) {
	tokenHandler := getParametrizedTokenHandler(validKey)
	fakeTokenHandler := getParametrizedTokenHandler(invalidKey)

	ba := auth.NewBearerAuthenticator(tokenHandler)

	var endpointHit bool

	router := gin.New()
	g := router.Group("/secure", ba.Authorize(security.Admin))
	g.GET("/data", func(context *gin.Context) {
		endpointHit = true

		up, err := auth.GetAuthenticatedUser(context)
		require.NoError(t, err)

		context.JSON(http.StatusOK, up)
	})

	cases := map[string]struct {
		loggedInUser         *security.UserPrincipal
		tokenHandler         *security.JwtHandler
		expectedStatusCode   int
		expectedErrorMessage string
	}{
		"success": {
			loggedInUser: &security.UserPrincipal{
				UserId: "12345678",
				Roles:  []security.Role{security.Student, security.Admin},
			},
			tokenHandler:         tokenHandler,
			expectedStatusCode:   http.StatusOK,
			expectedErrorMessage: "",
		},
		"fake_token": {
			loggedInUser: &security.UserPrincipal{
				UserId: "12345678",
				Roles:  []security.Role{security.Student, security.Admin},
			},
			tokenHandler:         fakeTokenHandler,
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: `{"message":"Unauthorized"}`,
		},
		"no_authorization_header": {
			loggedInUser:         nil,
			tokenHandler:         nil,
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: `{"message":"Unauthorized"}`,
		},
		"authorization_failure": {
			loggedInUser: &security.UserPrincipal{
				UserId: "12345678",
				Roles:  []security.Role{security.Student},
			},
			tokenHandler:         tokenHandler,
			expectedStatusCode:   http.StatusForbidden,
			expectedErrorMessage: `{"message":"Forbidden. Required user role: admin"}`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {

			endpointHit = false

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/secure/data", nil)
			if tc.loggedInUser != nil {
				tokenString, err := tc.tokenHandler.Generate(tc.loggedInUser)
				require.NoError(t, err)
				req.Header.Add("Authorization", "Bearer "+tokenString)
			}
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)

			if tc.expectedStatusCode == http.StatusOK {
				assert.True(t, endpointHit)
				var userPrincipal security.UserPrincipal
				err := json.Unmarshal(w.Body.Bytes(), &userPrincipal)
				assert.NoError(t, err)
				assert.Equal(t, tc.loggedInUser, &userPrincipal)
			} else {
				assert.False(t, endpointHit)
				assert.Equal(t, tc.expectedErrorMessage, w.Body.String())
			}
		})
	}
}

func TestBearerAuthenticator_Integration_Authenticate(t *testing.T) {
	tokenHandler := getParametrizedTokenHandler(validKey)
	fakeTokenHandler := getParametrizedTokenHandler(invalidKey)

	ba := auth.NewBearerAuthenticator(tokenHandler)

	var endpointHit bool

	router := gin.New()
	g := router.Group("/secure", ba.Authenticate)
	g.GET("/data", func(context *gin.Context) {
		endpointHit = true

		up, err := auth.GetAuthenticatedUser(context)
		require.NoError(t, err)

		context.JSON(http.StatusOK, up)
	})

	cases := map[string]struct {
		loggedInUser         *security.UserPrincipal
		tokenHandler         *security.JwtHandler
		expectedStatusCode   int
		expectedErrorMessage string
	}{
		"success_multiple_roles": {
			loggedInUser: &security.UserPrincipal{
				UserId: "12345678",
				Roles:  []security.Role{security.Student, security.Admin},
			},
			tokenHandler:         tokenHandler,
			expectedStatusCode:   http.StatusOK,
			expectedErrorMessage: "",
		},
		"fake_token": {
			loggedInUser: &security.UserPrincipal{
				UserId: "12345678",
				Roles:  []security.Role{security.Student, security.Admin},
			},
			tokenHandler:         fakeTokenHandler,
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: `{"message":"Unauthorized"}`,
		},
		"no_authorization_header": {
			loggedInUser:         nil,
			tokenHandler:         nil,
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: `{"message":"Unauthorized"}`,
		},
		"success_single_role": {
			loggedInUser: &security.UserPrincipal{
				UserId: "12345678",
				Roles:  []security.Role{security.Student},
			},
			tokenHandler:         tokenHandler,
			expectedStatusCode:   http.StatusOK,
			expectedErrorMessage: "",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {

			endpointHit = false

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/secure/data", nil)
			if tc.loggedInUser != nil {
				tokenString, err := tc.tokenHandler.Generate(tc.loggedInUser)
				require.NoError(t, err)
				req.Header.Add("Authorization", "Bearer "+tokenString)
			}
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)

			if tc.expectedStatusCode == http.StatusOK {
				assert.True(t, endpointHit)
				var userPrincipal security.UserPrincipal
				err := json.Unmarshal(w.Body.Bytes(), &userPrincipal)
				assert.NoError(t, err)
				assert.Equal(t, tc.loggedInUser, &userPrincipal)
			} else {
				assert.False(t, endpointHit)
				assert.Equal(t, tc.expectedErrorMessage, w.Body.String())
			}
		})
	}
}
