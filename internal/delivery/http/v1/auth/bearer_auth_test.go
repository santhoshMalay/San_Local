package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockHttp "github.com/zhuravlev-pe/course-watch/internal/delivery/http/mocks"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testSetup struct {
	router *gin.Engine
	bth    *mockHttp.MockBearerTokenHandler
	ba     *BearerAuthenticator
}

func getTestSetup(t *testing.T) *testSetup {
	t.Helper()
	ctrl := gomock.NewController(t)

	var ts testSetup

	ts.bth = mockHttp.NewMockBearerTokenHandler(ctrl)
	ts.ba = NewBearerAuthenticator(ts.bth)
	ts.router = gin.New()

	return &ts
}

const testData = "test data"
const validToken = "valid.token.value"
const invalidToken = "42"
const unauthorizedMessageBody = `{"message":"Unauthorized"}`

var referencePayload = &security.JwtPayload{
	UserPrincipal: security.UserPrincipal{
		UserId: "1111111",
		Roles:  []security.Role{security.Student},
	},
}

func TestBearerAuthenticator_Authenticate(t *testing.T) {
	ts := getTestSetup(t)

	var endpointHit bool

	g := ts.router.Group("/secure", ts.ba.Authenticate)
	g.GET("/data", func(context *gin.Context) {
		endpointHit = true
		context.String(http.StatusOK, testData)
	})

	cases := map[string]struct {
		header              string
		expectedParseInput  string
		expectedParseOutput *security.JwtPayload
		expectedParseError  error
		expectedStatusCode  int
		expectedBody        string
	}{
		"success": {
			header:              "Bearer " + validToken,
			expectedParseInput:  validToken,
			expectedParseOutput: referencePayload,
			expectedParseError:  nil,
			expectedStatusCode:  http.StatusOK,
			expectedBody:        testData,
		},
		"invalid_token": {
			header:              "Bearer " + invalidToken,
			expectedParseInput:  invalidToken,
			expectedParseOutput: nil,
			expectedParseError:  errors.New("some parsing error"),
			expectedStatusCode:  http.StatusUnauthorized,
			expectedBody:        unauthorizedMessageBody,
		},
		"missing_token": {
			header:              "Bearer ",
			expectedParseInput:  "",
			expectedParseOutput: nil,
			expectedParseError:  nil,
			expectedStatusCode:  http.StatusUnauthorized,
			expectedBody:        unauthorizedMessageBody,
		},
		"malformed_header": {
			header:              "Bearer more parts",
			expectedParseInput:  "",
			expectedParseOutput: nil,
			expectedParseError:  nil,
			expectedStatusCode:  http.StatusUnauthorized,
			expectedBody:        unauthorizedMessageBody,
		},
		"whitespace_header": {
			header:              "      ",
			expectedParseInput:  "",
			expectedParseOutput: nil,
			expectedParseError:  nil,
			expectedStatusCode:  http.StatusUnauthorized,
			expectedBody:        unauthorizedMessageBody,
		},
		"not_bearer": {
			header:              "Basic some_data",
			expectedParseInput:  "",
			expectedParseOutput: nil,
			expectedParseError:  nil,
			expectedStatusCode:  http.StatusUnauthorized,
			expectedBody:        unauthorizedMessageBody,
		},
		"no_header": {
			header:              "",
			expectedParseInput:  "",
			expectedParseOutput: nil,
			expectedParseError:  nil,
			expectedStatusCode:  http.StatusUnauthorized,
			expectedBody:        unauthorizedMessageBody,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if c.expectedParseInput != "" {
				ts.bth.EXPECT().Parse(c.expectedParseInput).Times(1).Return(c.expectedParseOutput, c.expectedParseError)
			}

			endpointHit = false

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/secure/data", nil)
			if c.header != "" {
				req.Header.Add("Authorization", c.header)
			}

			ts.router.ServeHTTP(w, req)

			require.Equal(t, c.expectedStatusCode, w.Code)
			require.Equal(t, c.expectedBody, w.Body.String())
			if c.expectedStatusCode == http.StatusOK {
				require.True(t, endpointHit)
			} else {
				require.False(t, endpointHit)
			}
		})
	}
}

func TestBearerAuthenticator_GenerateToken(t *testing.T) {
	ctrl := gomock.NewController(t)

	bth := mockHttp.NewMockBearerTokenHandler(ctrl)
	ba := NewBearerAuthenticator(bth)

	bth.EXPECT().Generate(&referencePayload.UserPrincipal).Times(1).Return(validToken, nil)

	token, err := ba.GenerateToken(&referencePayload.UserPrincipal)

	require.Equal(t, validToken, token)
	require.NoError(t, err)
}

func TestGetAuthenticatedUser_Success(t *testing.T) {
	ts := getTestSetup(t)

	var endpointHit bool

	g := ts.router.Group("/secure", ts.ba.Authenticate)
	g.GET("/data", func(context *gin.Context) {
		endpointHit = true
		up, err := GetAuthenticatedUser(context)
		require.Equal(t, &referencePayload.UserPrincipal, up)
		require.NoError(t, err)
		context.String(http.StatusOK, testData)
	})

	ts.bth.EXPECT().Parse(validToken).Times(1).Return(referencePayload, nil)

	endpointHit = false

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/secure/data", nil)
	req.Header.Add("Authorization", "Bearer "+validToken)

	ts.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, testData, w.Body.String())
	require.True(t, endpointHit)

}

func TestGetAuthenticatedUser_NoMiddleware(t *testing.T) {
	ts := getTestSetup(t)

	var endpointHit bool

	g := ts.router.Group("/secure")
	g.GET("/data", func(context *gin.Context) {
		endpointHit = true
		up, err := GetAuthenticatedUser(context)
		require.Nil(t, up)
		require.Error(t, err)
		context.String(http.StatusOK, testData)
	})

	ts.bth.EXPECT().Parse(validToken).Times(0)

	endpointHit = false

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/secure/data", nil)
	req.Header.Add("Authorization", "Bearer "+validToken)

	ts.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, testData, w.Body.String())
	require.True(t, endpointHit)
}
