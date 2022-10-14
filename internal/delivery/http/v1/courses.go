package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initCoursesRoutes(api *gin.RouterGroup) {
	courses := api.Group("/courses")
	{
		//courses.GET("", h.getAllCourses)
		courses.GET("/:id", h.getCourseById)
	}
}

// @Summary Get Course By course id
// @Tags courses
// @Description  get course by id
// @ModuleID getCourseById
// @Accept  json
// @Produce  json
// @Param id path string true "course id"
// @Success 200 {object} core.Course
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /courses/{id} [get]
func (h *Handler) getCourseById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "empty id param")
		return
	}

	course, err := h.services.Courses.GetById(c.Request.Context(), id)
	if err != nil {
		//TODO: 404 Not found
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, course)
}
