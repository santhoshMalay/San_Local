package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhuravlev-pe/course-watch/internal/service"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
)

type Handler struct {
	services *service.Services
	bearer   BearerAuthenticator
	//TODO: logger
}

type BearerAuthenticator interface {
	Authenticate(ctx *gin.Context)
	Authorize(role security.Role) func(ctx *gin.Context)
	GenerateToken(principal *security.UserPrincipal) (string, error)
}

func NewHandler(services *service.Services, bearer BearerAuthenticator) *Handler {
	return &Handler{
		services: services,
		bearer:   bearer,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initCoursesRoutes(v1)
		h.initUserRoutes(v1)
	}
}
