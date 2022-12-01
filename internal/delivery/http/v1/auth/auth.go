package auth

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_auth.go

import (
	"github.com/zhuravlev-pe/course-watch/pkg/security"
	"time"
)

type BearerTokenHandler interface {
	Generate(principal *security.UserPrincipal) (string, error)
	Parse(tokenString string) (*security.JwtPayload, error)
	GetTokenTtl() time.Duration
}
