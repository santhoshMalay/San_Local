package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/zhuravlev-pe/course-watch/internal/delivery/http/v1/utils"
	"github.com/zhuravlev-pe/course-watch/internal/repository"
	"net/http"
)

func (h *Handler) parseRequestBody(ctx *gin.Context, input interface{}) bool {
	if err := ctx.BindJSON(input); err != nil {
		utils.ErrorResponseMessageOverride(ctx, http.StatusBadRequest, err, "body is missing or invalid")
		return false
	}
	return true
}

func (h *Handler) handleServiceError(ctx *gin.Context, err error) {
	if err == repository.ErrNotFound {
		utils.ErrorResponse(ctx, http.StatusNotFound, err)
		return
	}

	var validationErrors validation.Errors
	if errors.As(err, &validationErrors) {
		utils.ValidationErrorResponse(ctx, validationErrors)
		return
	}

	utils.ErrorResponseMessageOverride(ctx, http.StatusInternalServerError, err, "internal server error")
	return
}
