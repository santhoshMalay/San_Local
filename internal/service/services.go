package service

import (
	"context"
	"time"

	"github.com/zhuravlev-pe/course-watch/internal/core"
	"github.com/zhuravlev-pe/course-watch/internal/repository"
	"github.com/zhuravlev-pe/course-watch/pkg/idgen"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
)

type CreateCourseInput struct {
	Title       string
	Description string
}

type Courses interface {
	GetById(ctx context.Context, id string) (*core.Course, error)
	Create(ctx context.Context, input CreateCourseInput) (*core.Course, error)
}

type GetUserInfoOutput struct {
	Id               string          `json:"id"`
	Email            string          `json:"email"`
	FirstName        string          `json:"first_name"`
	LastName         string          `json:"last_name"`
	DisplayName      string          `json:"display_name"`
	RegistrationDate time.Time       `json:"registration_date"`
	Roles            []security.Role `json:"roles"`
}

type UpdateUserInfoInput struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DisplayName string `json:"display_name"`
}

type LoginInput struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Persistent bool   `json:"persistent"`
}

type SignupInput struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DisplayName string `json:"display_name"`
}

type PostUserLoginOutput struct {
	UserId      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type Users interface {
	GetUserInfo(ctx context.Context, id string) (*GetUserInfoOutput, error)
	UpdateUserInfo(ctx context.Context, id string, input *UpdateUserInfoInput) error
	Login(ctx context.Context, input *LoginInput) (*core.User, error)
	Signup(ctx context.Context, input *SignupInput) error
}

type Services struct {
	Courses Courses
	Users   Users
}

type Deps struct {
	Repos *repository.Repositories
	IdGen *idgen.IdGen
}

func NewServices(deps Deps) *Services {
	coursesService := NewCoursesService(deps.Repos.Courses, deps.IdGen)
	usersSrv := newUsersService(deps.Repos.Users, deps.IdGen)

	return &Services{
		Courses: coursesService,
		Users:   usersSrv,
	}
}
