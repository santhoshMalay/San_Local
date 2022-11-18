package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhuravlev-pe/course-watch/internal/delivery/http/v1/auth"
	"github.com/zhuravlev-pe/course-watch/internal/delivery/http/v1/utils"
	"github.com/zhuravlev-pe/course-watch/internal/repository"
	"github.com/zhuravlev-pe/course-watch/internal/service"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
)

func (h *Handler) initUserRoutes(api *gin.RouterGroup) {
	courses := api.Group("/user", h.bearer.Authenticate)
	{
		courses.GET("/", h.getUserInfo)
		courses.PUT("/", h.updateUserInfo)
	}
}

// @Summary Retrieve current user data
// @Tags User
// @Description returns info on the currently logged-in user. User_id is extracted from the bearer token
// @ModuleID getUserInfo
// @Accept  json
// @Produce  json
// @Success 200 {object} service.GetUserInfoOutput
// @Failure 400,401,404,500 {object} utils.Response
// @Failure default {object} utils.Response
// @Router /user/ [get]
func (h *Handler) getUserInfo(ctx *gin.Context) {
	up, err := auth.GetAuthenticatedUser(ctx)
	if err != nil {
		err = fmt.Errorf("authentication middleware failure: %w", err)
		utils.ErrorResponseMessageOverride(ctx, http.StatusInternalServerError, err, "user data processing failure")
		return
	}

	result, err := h.services.Users.GetUserInfo(ctx.Request.Context(), up.UserId)

	if err != nil {
		// TODO: discriminate between validation errors, logic errors and internal server errors
		if err == repository.ErrNotFound {
			utils.ErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// @Summary Modify current user data
// @Tags User
// @Description modifies user info for the currently logged-in user. User_id is extracted from the bearer token
// @ModuleID updateUserInfo
// @Accept  json
// @Produce  json
// @Param input body service.UpdateUserInfoInput true "user info"
// @Success 200
// @Failure 400,401,404,500 {object} utils.Response
// @Failure default {object} utils.Response
// @Router /user/ [put]
func (h *Handler) updateUserInfo(ctx *gin.Context) {
	up, err := auth.GetAuthenticatedUser(ctx)
	if err != nil {
		err = fmt.Errorf("authentication middleware failure: %w", err)
		utils.ErrorResponseMessageOverride(ctx, http.StatusInternalServerError, err, "user data processing failure")
		return
	}
	var input service.UpdateUserInfoInput
	if err = ctx.BindJSON(&input); err != nil {
		utils.ErrorResponseString(ctx, http.StatusBadRequest, "invalid input body")
		return
	}

	err = h.services.Users.UpdateUserInfo(ctx.Request.Context(), up.UserId, &input)

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

func (h *Handler) initAuthRoutes(api *gin.RouterGroup) {
	courses := api.Group("/auth")
	{
		courses.POST("/login", h.userLogin)
		courses.POST("/signup", h.signupNewUser)
	}
}

// @Summary Authenticate user credentials
// @Tags Authentication
// @Description authenticates the user log-in credentials
// @ModuleID userLogin
// @Accept  json
// @Produce  json
// @Success 200 {object} service.LoginInput
// @Failure 400,500 {object} utils.Response
// @Router /auth [Post]
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

// @Summary New user signup
// @Tags Authentication
// @Description Creates new user with the given detials
// @ModuleID signupNewUser
// @Accept  json
// @Produce  json
// @Success 200 {object} service.LoginInput
// @Failure 400,500 {object} utils.Response
// @Router /auth [Post]
func (h *Handler) signupNewUser(ctx *gin.Context) {
	var input service.SignupInput
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

	// up := security.UserPrincipal{UserId: result.Id, Roles: result.Roles}
	// token, err := h.bearer.GenerateToken(&up)
	// if err != nil {
	// 	utils.ErrorResponse(ctx, http.StatusInternalServerError, err)
	// 	return
	// }
	// output := service.PostUserLoginOutput{
	// 	UserId:      up.UserId,
	// 	AccessToken: token,
	// 	ExpiresIn:   0,
	// }
	ctx.Status(http.StatusOK)
}
