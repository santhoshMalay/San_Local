package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhuravlev-pe/course-watch/internal/delivery/http/v1/auth"
	"github.com/zhuravlev-pe/course-watch/internal/service"
)

type Handler struct {
	services *service.Services
	bearer   auth.BearerTokenHandler
	//TODO: logger
}

func NewHandler(services *service.Services, jwtHandler auth.BearerTokenHandler) *Handler {
	return &Handler{
		services: services,
		bearer:   jwtHandler,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initCoursesRoutes(v1)
	}
}
