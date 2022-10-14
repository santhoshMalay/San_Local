package service

import (
	"context"
	"github.com/zhuravlev-pe/course-watch/internal/core"
	"github.com/zhuravlev-pe/course-watch/internal/repository"
)

type CoursesService struct {
	repo repository.Courses
}

func NewCoursesService(repo repository.Courses) *CoursesService {
	return &CoursesService{
		repo: repo,
	}
}

func (s *CoursesService) GetById(ctx context.Context, id string) (*core.Course, error) {
	return s.repo.GetById(ctx, id)
}
