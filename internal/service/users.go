package service

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/zhuravlev-pe/course-watch/internal/core"
	"github.com/zhuravlev-pe/course-watch/internal/repository"
	"github.com/zhuravlev-pe/course-watch/pkg/idgen"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
	"golang.org/x/crypto/bcrypt"
)

type usersService struct {
	repo  repository.Users
	idGen *idgen.IdGen
}

func (u *usersService) GetUserInfo(ctx context.Context, id string) (*GetUserInfoOutput, error) {
	user, err := u.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	var result GetUserInfoOutput
	result.Id = user.Id
	result.Email = user.Email
	result.FirstName = user.FirstName
	result.LastName = user.LastName
	result.DisplayName = user.DisplayName
	result.RegistrationDate = user.RegistrationDate
	result.Roles = user.Roles
	return &result, nil
}

func (i *UpdateUserInfoInput) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.FirstName, validation.Required),
		validation.Field(&i.LastName, validation.Required),
	)
}

func (u *usersService) UpdateUserInfo(ctx context.Context, id string, input *UpdateUserInfoInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	_, err := u.repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	var upd repository.UpdateUserInput
	upd.FirstName = input.FirstName
	upd.LastName = input.LastName
	upd.DisplayName = input.DisplayName
	return u.repo.Update(ctx, id, &upd)
}

func (i *SignupUserInput) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.Email, validation.Required, is.Email),
		validation.Field(&i.Password, validation.Required, validation.Length(8, 20)),
		validation.Field(&i.FirstName, validation.Required),
		validation.Field(&i.LastName, validation.Required),
	)
}

func (u *usersService) Signup(ctx context.Context, input *SignupUserInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	user, err := u.repo.GetByEmail(ctx, input.Email)
	if err != nil && err != repository.ErrNotFound {
		//Any error other than ErrorNotFound should stop the Signup flow as ErrorNotFound is valid for the user Signup
		return err
	}

	if user != nil && user.Email == input.Email {
		return ErrUserAlreadyExist
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user = &core.User{
		Id:               u.idGen.Generate(),
		Email:            input.Email,
		FirstName:        input.FirstName,
		LastName:         input.LastName,
		DisplayName:      input.DisplayName,
		RegistrationDate: time.Now(),
		HashedPassword:   hashPassword,
		Roles:            []security.Role{security.Student},
	}

	if err = u.repo.Insert(ctx, user); err != nil {
		return err
	}
	return nil
}

func (u *usersService) Login(ctx context.Context, input *LoginInput) (*core.User, error) {
	user, err := u.repo.GetByEmail(ctx, input.Email)
	if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	if err == repository.ErrNotFound {
		return nil, ErrInvalidCredentials
	}

	if err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}

func newUsersService(repo repository.Users, idGen *idgen.IdGen) Users {
	return &usersService{
		repo:  repo,
		idGen: idGen,
	}
}
