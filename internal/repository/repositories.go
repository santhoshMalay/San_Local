package repository

import (
	"context"
	"github.com/zhuravlev-pe/course-watch/internal/core"
)

type Courses interface {
	GetById(ctx context.Context, id string) (*core.Course, error)
}

type Repositories struct {
	Courses Courses
}
