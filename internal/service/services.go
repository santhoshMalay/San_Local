package service

import (
	"context"
	"github.com/zhuravlev-pe/course-watch/internal/core"
	"github.com/zhuravlev-pe/course-watch/internal/repository"
)

type Courses interface {
	GetById(ctx context.Context, id string) (*core.Course, error)
}

type Services struct {
	Courses Courses
}

type Deps struct {
	Repos *repository.Repositories
}

func NewServices(deps Deps) *Services {
	coursesService := NewCoursesService(deps.Repos.Courses)

	return &Services{
		Courses: coursesService,
	}
}
