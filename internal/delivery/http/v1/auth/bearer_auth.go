package auth

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "github.com/zhuravlev-pe/course-watch/internal/delivery/http/v1/utils"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
	"net/http"
	"strings"
	"time"
)

const userKey = "user_principal"

type BearerAuthenticator struct {
	tokenHandler BearerTokenHandler
}

func NewBearerAuthenticator(tokenHandler BearerTokenHandler) *BearerAuthenticator {
	return &BearerAuthenticator{tokenHandler: tokenHandler}
}

// Authenticate implements authentication middleware. When used on a group router, child endpoints will be called only
// if a valid bearer token is passed in the request. These endpoints may call GetAuthenticatedUser() to access user
// data
//
// Warning: if applied to an endpoint, the auth failure will abort the middleware chain but WILL NOT prevent the
// endpoint handler from running, due to the way how gin works:
// https://github.com/gin-gonic/gin/issues/2442
// Check gin.Context.IsAborted() in the handler code to ensure that authentication has passed
func (ba *BearerAuthenticator) Authenticate(ctx *gin.Context) {
	payload, err := ba.parseAuthHeader(ctx)
	if err != nil {
		// We do not want to report the error details to the caller here in order to avoid revealing security related
		// info. ErrorResponseWithMessage() should log the error
		v1.ErrorResponseMessageOverride(ctx, http.StatusUnauthorized, err, "Unauthorized")
		return
	}
	up := payload.UserPrincipal
	ctx.Set(userKey, &up)
}

func (ba *BearerAuthenticator) parseAuthHeader(ctx *gin.Context) (*security.JwtPayload, error) {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		return nil, errors.New("empty or missing 'Authorization' header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("invalid 'Authorization' header")
	}

	if len(headerParts[1]) == 0 {
		return nil, errors.New("bearer token is empty")
	}

	return ba.tokenHandler.Parse(headerParts[1])
}

// Authorize middleware performs authentication and ensures that user has the specified role
//
// Waring: same caveat as for the Authenticate() middleware. Apply to group middleware only
func (ba *BearerAuthenticator) Authorize(role security.Role) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ba.Authenticate(ctx)
		if ctx.IsAborted() {
			return
		}
		EnsureAuthorizedUser(ctx, role)
	}
}

// EnsureAuthorizedUser checks if user is authenticated and has a required role. If not, aborts the context with 403
// (Forbidden) and a proper message and returns false. Handler method should call return in this case. Returns true
// otherwise
//
//		func myEndpoint(ctx *gin.Context) {
//		    if !EnsureAuthorizedUser(ctx, security.Admin) {
//		        return
//		    }
//	     // ... do admin stuff ...
//		}
func EnsureAuthorizedUser(ctx *gin.Context, role security.Role) bool {
	up, err := GetAuthenticatedUser(ctx)
	if err != nil {
		// No authentication middleware
		msg := getUnauthorizedMessage(role)
		v1.ErrorResponseMessageOverride(ctx, http.StatusForbidden, err, msg)
		return false
	}
	if !up.HasRole(role) {
		// User does not have the required role
		msg := getUnauthorizedMessage(role)
		v1.ErrorResponseString(ctx, http.StatusForbidden, msg)
		return false
	}
	return true
}

func getUnauthorizedMessage(role security.Role) string {
	return fmt.Sprintf("Forbidden. Required user role: %s", role.String())
}

// GenerateToken creates a signed token string to be passed to the frontend in a response
func (ba *BearerAuthenticator) GenerateToken(principal *security.UserPrincipal) (string, error) {
	return ba.tokenHandler.Generate(principal)
}

func (ba *BearerAuthenticator) GetTokenTtl() time.Duration {
	return ba.tokenHandler.GetTokenTtl()
}

// GetAuthenticatedUser returns the authenticated user data when called in endpoints protected by the
// BearerAuthenticator.Authenticate middleware
func GetAuthenticatedUser(ctx *gin.Context) (*security.UserPrincipal, error) {
	data, ok := ctx.Get(userKey)
	if !ok {
		return nil, errors.New("user data is missing. Ensure that authentication middleware is properly called before accessing the user data")
	}
	up, ok := data.(*security.UserPrincipal)
	if !ok || up == nil {
		return nil, errors.New("user data is empty. Ensure that authentication middleware is properly called before accessing the user data")
	}
	return up, nil
}
