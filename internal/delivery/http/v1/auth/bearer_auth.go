package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	deliveryHttp "github.com/zhuravlev-pe/course-watch/internal/delivery/http"
	v1 "github.com/zhuravlev-pe/course-watch/internal/delivery/http/v1"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
	"net/http"
	"strings"
)

const userKey = "user_principal"

type BearerAuthenticator struct {
	tokenHandler deliveryHttp.BearerTokenHandler
}

func NewBearerAuthenticator(tokenHandler deliveryHttp.BearerTokenHandler) *BearerAuthenticator {
	return &BearerAuthenticator{tokenHandler: tokenHandler}
}

// Authenticate implements authentication middleware
func (ba *BearerAuthenticator) Authenticate(c *gin.Context) {
	payload, err := ba.parseAuthHeader(c)
	if err != nil {
		// We do not want to report the error details to the caller here in order to avoid revealing security related
		// info. ErrorResponseWithMessage() should log the error
		v1.ErrorResponseMessageOverride(c, http.StatusUnauthorized, err, "Unauthorized")
		return
	}
	up := payload.UserPrincipal
	c.Set(userKey, &up)
}

func (ba *BearerAuthenticator) parseAuthHeader(c *gin.Context) (*security.JwtPayload, error) {
	header := c.GetHeader("Authorization")
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

// GenerateToken creates a signed token string to be passed to the frontend in a response
func (ba *BearerAuthenticator) GenerateToken(principal *security.UserPrincipal) (string, error) {
	return ba.tokenHandler.Generate(principal)
}

func GetAuthenticatedUser(c *gin.Context) (*security.UserPrincipal, error) {
	data, ok := c.Get(userKey)
	if !ok {
		return nil, errors.New("user data is missing. Ensure that authentication middleware is properly called before accessing the user data")
	}
	up, ok := data.(*security.UserPrincipal)
	if !ok || up == nil {
		return nil, errors.New("user data is empty. Ensure that authentication middleware is properly called before accessing the user data")
	}
	return up, nil
}
