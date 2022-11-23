package v1

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zhuravlev-pe/course-watch/internal/delivery/http/v1/auth"
	"github.com/zhuravlev-pe/course-watch/internal/repository"
	"github.com/zhuravlev-pe/course-watch/internal/service"
	serviceMocks "github.com/zhuravlev-pe/course-watch/internal/service/mocks"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	iss      = "course-watch"
	aud      = "course-watch-api"
	tokenTtl = time.Hour * 1
)

var validKey = []byte("1234")

var sampleUserInfo = &service.GetUserInfoOutput{
	Id:               "1582550893222432768",
	Email:            "doe.j@example.com",
	FirstName:        "John",
	LastName:         "Doe",
	DisplayName:      "JonnyD",
	RegistrationDate: time.Date(2017, time.July, 21, 17, 32, 28, 0, time.UTC),
	Roles:            []security.Role{security.Student},
}

var sampleUserPrincipal = &security.UserPrincipal{
	UserId: "1582550893222432768",
	Roles:  []security.Role{security.Student},
}

type testSetup struct {
	router          *gin.Engine
	users           *serviceMocks.MockUsers
	handler         *Handler
	sampleUserToken string
}

func getTestSetup(t *testing.T) *testSetup {
	t.Helper()
	mockCtrl := gomock.NewController(t)
	mockUsers := serviceMocks.NewMockUsers(mockCtrl)
	var s service.Services
	s.Users = mockUsers

	jwt := security.NewJwtHandler(iss, aud, []string{aud}, tokenTtl, validKey)
	bearer := auth.NewBearerAuthenticator(jwt)
	token, err := bearer.GenerateToken(sampleUserPrincipal)
	require.NoError(t, err)

	handler := NewHandler(&s, bearer)

	router := gin.New()
	handler.Init(router.Group("/api"))

	return &testSetup{
		router:          router,
		users:           mockUsers,
		handler:         handler,
		sampleUserToken: token,
	}
}

var someDatabaseError = errors.New("some database error")

func addAuthorizationHeader(request *http.Request, setup *testSetup) {
	request.Header.Add("Authorization", "Bearer "+setup.sampleUserToken)
}

func TestGetUserInfo(t *testing.T) {
	cases := map[string]struct {
		setupMocks     func(ctx context.Context, mockUsers *serviceMocks.MockUsers)
		prepareRequest func(request *http.Request, setup *testSetup)
		responseCode   int
		responseBody   string
	}{
		"success": {
			setupMocks: func(ctx context.Context, mockUsers *serviceMocks.MockUsers) {
				mockUsers.EXPECT().GetUserInfo(ctx, sampleUserPrincipal.UserId).Return(sampleUserInfo, nil).Times(1)
			},
			prepareRequest: addAuthorizationHeader,
			responseCode:   http.StatusOK,
			responseBody:   `{"id":"1582550893222432768","email":"doe.j@example.com","first_name":"John","last_name":"Doe","display_name":"JonnyD","registration_date":"2017-07-21T17:32:28Z","roles":["student"]}`,
		},
		"not_found": {
			setupMocks: func(ctx context.Context, mockUsers *serviceMocks.MockUsers) {
				mockUsers.EXPECT().GetUserInfo(ctx, sampleUserPrincipal.UserId).Return(nil, repository.ErrNotFound).Times(1)
			},
			prepareRequest: addAuthorizationHeader,
			responseCode:   http.StatusNotFound,
			responseBody:   `{"title":"not found","status":404}`,
		},
		"internal_server_err": {
			setupMocks: func(ctx context.Context, mockUsers *serviceMocks.MockUsers) {
				mockUsers.EXPECT().GetUserInfo(ctx, sampleUserPrincipal.UserId).Return(nil, someDatabaseError).Times(1)
			},
			prepareRequest: addAuthorizationHeader,
			responseCode:   http.StatusInternalServerError,
			responseBody:   `{"title":"internal server error","status":500}`,
		},
		"unauthorized": {
			setupMocks:     func(ctx context.Context, mockUsers *serviceMocks.MockUsers) {},
			prepareRequest: func(request *http.Request, setup *testSetup) {},
			responseCode:   http.StatusUnauthorized,
			responseBody:   `{"title":"Unauthorized","status":401}`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			setup := getTestSetup(t)
			ctx := context.Background()
			tc.setupMocks(ctx, setup.users)

			request := httptest.NewRequest(http.MethodGet, "/api/v1/user", nil)
			tc.prepareRequest(request, setup)
			rec := httptest.NewRecorder()

			setup.router.ServeHTTP(rec, request)

			assert.Equal(t, tc.responseCode, rec.Code)
			assert.Equal(t, tc.responseBody, rec.Body.String())
		})
	}
}

func TestUpdateUserInfo(t *testing.T) {
	cases := map[string]struct {
		requestBody    string
		setupMocks     func(ctx context.Context, mockUsers *serviceMocks.MockUsers)
		prepareRequest func(request *http.Request, setup *testSetup)
		responseCode   int
		responseBody   string
	}{
		"success": {
			requestBody: `{"first_name":"UpdatedFirstName","last_name":"UpdatedLastName","display_name":"UpdatedDisplayName"}`,
			setupMocks: func(ctx context.Context, mockUsers *serviceMocks.MockUsers) {
				input := &service.UpdateUserInfoInput{
					FirstName:   "UpdatedFirstName",
					LastName:    "UpdatedLastName",
					DisplayName: "UpdatedDisplayName",
				}
				mockUsers.EXPECT().UpdateUserInfo(ctx, sampleUserPrincipal.UserId, input).Return(nil).Times(1)
			},
			prepareRequest: addAuthorizationHeader,
			responseCode:   http.StatusNoContent,
			responseBody:   "",
		},
		"empty_body": {
			requestBody:    "",
			setupMocks:     func(ctx context.Context, mockUsers *serviceMocks.MockUsers) {},
			prepareRequest: addAuthorizationHeader,
			responseCode:   http.StatusBadRequest,
			responseBody:   `{"title":"body is missing or invalid","status":400}`,
		},
		"invalid_json": {
			requestBody:    "not a valid json",
			setupMocks:     func(ctx context.Context, mockUsers *serviceMocks.MockUsers) {},
			prepareRequest: addAuthorizationHeader,
			responseCode:   http.StatusBadRequest,
			responseBody:   `{"title":"body is missing or invalid","status":400}`,
		},
		"not_found": {
			requestBody: `{"first_name":"UpdatedFirstName","last_name":"UpdatedLastName","display_name":"UpdatedDisplayName"}`,
			setupMocks: func(ctx context.Context, mockUsers *serviceMocks.MockUsers) {
				input := &service.UpdateUserInfoInput{
					FirstName:   "UpdatedFirstName",
					LastName:    "UpdatedLastName",
					DisplayName: "UpdatedDisplayName",
				}
				mockUsers.EXPECT().UpdateUserInfo(ctx, sampleUserPrincipal.UserId, input).Return(repository.ErrNotFound).Times(1)
			},
			prepareRequest: addAuthorizationHeader,
			responseCode:   http.StatusNotFound,
			responseBody:   `{"title":"not found","status":404}`,
		},
		"internal_server_err": {
			requestBody: `{"first_name":"UpdatedFirstName","last_name":"UpdatedLastName","display_name":"UpdatedDisplayName"}`,
			setupMocks: func(ctx context.Context, mockUsers *serviceMocks.MockUsers) {
				input := &service.UpdateUserInfoInput{
					FirstName:   "UpdatedFirstName",
					LastName:    "UpdatedLastName",
					DisplayName: "UpdatedDisplayName",
				}
				mockUsers.EXPECT().UpdateUserInfo(ctx, sampleUserPrincipal.UserId, input).Return(someDatabaseError).Times(1)
			},
			prepareRequest: addAuthorizationHeader,
			responseCode:   http.StatusInternalServerError,
			responseBody:   `{"title":"internal server error","status":500}`,
		},
		"unauthorized": {
			requestBody:    `{"first_name":"UpdatedFirstName","last_name":"UpdatedLastName","display_name":"UpdatedDisplayName"}`,
			setupMocks:     func(ctx context.Context, mockUsers *serviceMocks.MockUsers) {},
			prepareRequest: func(request *http.Request, setup *testSetup) {},
			responseCode:   http.StatusUnauthorized,
			responseBody:   `{"title":"Unauthorized","status":401}`,
		},
		"validation_failure": {
			requestBody: `{"last_name":"UpdatedLastName","display_name":"UpdatedDisplayName"}`,
			setupMocks: func(ctx context.Context, mockUsers *serviceMocks.MockUsers) {
				input := &service.UpdateUserInfoInput{
					LastName:    "UpdatedLastName",
					DisplayName: "UpdatedDisplayName",
				}
				validationError := input.Validate()
				mockUsers.EXPECT().UpdateUserInfo(ctx, sampleUserPrincipal.UserId, input).Return(validationError).Times(1)
			},
			prepareRequest: addAuthorizationHeader,
			responseCode:   http.StatusBadRequest,
			responseBody:   `{"title":"invalid request parameters","status":400,"validation_errors":{"first_name":"cannot be blank"}}`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			setup := getTestSetup(t)
			ctx := context.Background()
			tc.setupMocks(ctx, setup.users)
			var body io.Reader
			if tc.requestBody != "" {
				body = strings.NewReader(tc.requestBody)
			}

			request := httptest.NewRequest(http.MethodPut, "/api/v1/user", body)
			tc.prepareRequest(request, setup)
			rec := httptest.NewRecorder()

			setup.router.ServeHTTP(rec, request)

			assert.Equal(t, tc.responseCode, rec.Code)
			assert.Equal(t, tc.responseBody, rec.Body.String())
		})
	}
}
