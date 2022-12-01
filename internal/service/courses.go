package service

import (
	"context"

	"github.com/zhuravlev-pe/course-watch/internal/core"
	"github.com/zhuravlev-pe/course-watch/internal/repository"
	"github.com/zhuravlev-pe/course-watch/pkg/idgen"
)

type CoursesService struct {
	repo  repository.Courses
	idGen *idgen.IdGen
}

func NewCoursesService(repo repository.Courses, idGen *idgen.IdGen) *CoursesService {
	return &CoursesService{
		repo:  repo,
		idGen: idGen,
	}
}

func (s *CoursesService) GetById(ctx context.Context, id string) (*core.Course, error) {
	return s.repo.GetById(ctx, id)
}

func (s *CoursesService) Create(ctx context.Context, input CreateCourseInput) (*core.Course, error) {
	course := &core.Course{
		Id:          s.idGen.Generate(),
		Title:       input.Title,
		Description: input.Description,
	}
	if err := s.repo.Insert(ctx, course); err != nil {
		return nil, err
	}
	return course, nil
}
