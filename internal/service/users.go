package service

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/zhuravlev-pe/course-watch/internal/repository"
	"github.com/zhuravlev-pe/course-watch/pkg/idgen"
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

func newUsersService(repo repository.Users, idGen *idgen.IdGen) Users {
	return &usersService{
		repo:  repo,
		idGen: idGen,
	}
}
