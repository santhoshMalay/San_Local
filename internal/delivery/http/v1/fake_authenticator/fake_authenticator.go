package fake_authenticator

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/zhuravlev-pe/course-watch/internal/delivery/http/v1"
	"github.com/zhuravlev-pe/course-watch/internal/repository/fake_repo"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
)

type fakeBearerAuthenticator struct{}

func New() v1.BearerAuthenticator {
	return fakeBearerAuthenticator{}
}

func (f fakeBearerAuthenticator) Authenticate(ctx *gin.Context) {
	var up security.UserPrincipal
	up.UserId = fake_repo.SampleUser.Id
	up.Roles = fake_repo.SampleUser.Roles
	ctx.Set("user_principal", &up)
}

func (f fakeBearerAuthenticator) Authorize(_ security.Role) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		f.Authenticate(ctx)
	}
}

func (f fakeBearerAuthenticator) GenerateToken(_ *security.UserPrincipal) (string, error) {
	return "fake.bearer.token", nil
}
