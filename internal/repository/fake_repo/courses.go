package fake_repo

import (
	"context"

	"github.com/zhuravlev-pe/course-watch/internal/core"
	"github.com/zhuravlev-pe/course-watch/internal/repository"
)

type courses struct {
	data map[string]*core.Course
}

func NewCourses() repository.Courses {
	return &courses{
		data: map[string]*core.Course{},
	}
}

func (c *courses) GetById(_ context.Context, id string) (*core.Course, error) {
	course, ok := c.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return course, nil
}

func (c *courses) Insert(_ context.Context, course *core.Course) error {
	c.data[course.Id] = course
	return nil
}
