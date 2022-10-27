package http

import "github.com/zhuravlev-pe/course-watch/pkg/security"

// BearerTokenHandler represents a service for generating and validating security tokens used for bearer authentication
type BearerTokenHandler interface {
	// Generate creates a new signed token
	Generate(principal *security.UserPrincipal) (string, error)
	// Parse parses and validates a JWT token string
	Parse(tokenString string) (*security.JwtPayload, error)
}
