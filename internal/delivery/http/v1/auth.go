package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhuravlev-pe/course-watch/internal/delivery/http/v1/utils"
	"github.com/zhuravlev-pe/course-watch/internal/repository"
	"github.com/zhuravlev-pe/course-watch/internal/service"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
)

func (h *Handler) initAuthRoutes(api *gin.RouterGroup) {
	courses := api.Group("/auth")
	{
		courses.POST("/signup", h.signupNewUser)
		courses.POST("/login", h.userLogin)
	}
}

// @Summary New user signup
// @Tags Authentication
// @Description Creates new user with the given detials
// @ModuleID signupNewUser
// @Accept  json
// @Produce  json
// @Param input body service.SignupUserInput true "New user signup details"
// @Success 200 {object} service.LoginInput
// @Failure 400,500 {object} utils.Response
// @Router /auth/signup [Post]
func (h *Handler) signupNewUser(ctx *gin.Context) {
	var input service.SignupUserInput
	if err := ctx.BindJSON(&input); err != nil {
		utils.ErrorResponseString(ctx, http.StatusBadRequest, "invalid input body")
		return
	}
	err := h.services.Users.Signup(ctx.Request.Context(), &input)

	if err != nil {
		// TODO: discriminate between validation errors, logic errors and internal server errors
		if err == repository.ErrNotFound {
			utils.ErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Authenticate user credentials
// @Tags Authentication
// @Description authenticates the user log-in credentials
// @ModuleID userLogin
// @Accept  json
// @Produce  json
// @Param input body service.LoginInput true "Login user details"
// @Success 200 {object} service.LoginInput
// @Failure 400 {object} utils.Response
// @Router /auth/login [Post]
func (h *Handler) userLogin(ctx *gin.Context) {
	var input service.LoginInput
	if err := ctx.BindJSON(&input); err != nil {
		utils.ErrorResponseString(ctx, http.StatusBadRequest, "invalid input body")
		return
	}
	result, err := h.services.Users.Login(ctx.Request.Context(), &input)

	if err != nil {
		// TODO: discriminate between validation errors, logic errors and internal server errors
		if err == repository.ErrNotFound {
			utils.ErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	up := security.UserPrincipal{UserId: result.Id, Roles: result.Roles}
	token, err := h.bearer.GenerateToken(&up)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	output := service.PostUserLoginOutput{
		UserId:      up.UserId,
		AccessToken: token,
		ExpiresIn:   0,
	}
	ctx.JSON(http.StatusOK, output)
}
