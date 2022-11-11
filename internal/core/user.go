package core

import (
	"github.com/zhuravlev-pe/course-watch/pkg/security"
	"time"
)

type User struct {
	Id               string
	Email            string
	FirstName        string
	LastName         string
	DisplayName      string
	RegistrationDate time.Time
	HashedPassword   []byte
	Roles            []security.Role
}
