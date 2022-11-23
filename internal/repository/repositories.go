package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/repositories.go

import (
	"context"

	"github.com/zhuravlev-pe/course-watch/internal/core"
)

type Courses interface {
	GetById(ctx context.Context, id string) (*core.Course, error)
	Insert(ctx context.Context, course *core.Course) error
}

type UpdateUserInput struct {
	FirstName   string
	LastName    string
	DisplayName string
}

type Users interface {
	GetById(ctx context.Context, id string) (*core.User, error)
	Insert(ctx context.Context, user *core.User) error
	Update(ctx context.Context, id string, input *UpdateUserInput) error
	GetByEmail(ctx context.Context, email string) (*core.User, error)
}

type Repositories struct {
	Courses Courses
	Users   Users
}
