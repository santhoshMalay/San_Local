package service

import (
	"context"
	"github.com/zhuravlev-pe/course-watch/internal/core"
	"github.com/zhuravlev-pe/course-watch/internal/repository"
	"github.com/zhuravlev-pe/course-watch/pkg/idgen"
)

type CreateCourseInput struct {
	Title       string
	Description string
}

type Courses interface {
	GetById(ctx context.Context, id string) (*core.Course, error)
	Create(ctx context.Context, input CreateCourseInput) (*core.Course, error)
}

type Services struct {
	Courses Courses
}

type Deps struct {
	Repos *repository.Repositories
	IdGen *idgen.IdGen
}

func NewServices(deps Deps) *Services {
	coursesService := NewCoursesService(deps.Repos.Courses, deps.IdGen)

	return &Services{
		Courses: coursesService,
	}
}
