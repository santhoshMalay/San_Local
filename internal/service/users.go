package service

import (
	"context"

	"github.com/zhuravlev-pe/course-watch/internal/core"
	"github.com/zhuravlev-pe/course-watch/internal/repository"
	"github.com/zhuravlev-pe/course-watch/pkg/idgen"
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

func (u *usersService) UpdateUserInfo(ctx context.Context, id string, input *UpdateUserInfoInput) error {
	//TODO input validation
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

func (u *usersService) Signup(ctx context.Context, input *SignupUserInput) error {
	user, err := u.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return err
	}

	if user.Email == input.Email {
		return ErrUserAlreadyExist
	}

	user = &core.User{
		Id:          u.idGen.Generate(),
		Email:       input.Email,
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		DisplayName: input.DisplayName,
	}

	if err := u.repo.Insert(ctx, user); err != nil {
		return err
	}
	return nil
}

func (u *usersService) Login(ctx context.Context, input *LoginInput) (*core.User, error) {
	user, err := u.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(input.Password)); err != nil {
		return nil, err
	}
	return user, nil
}

func newUsersService(repo repository.Users, idGen *idgen.IdGen) Users {
	return &usersService{
		repo:  repo,
		idGen: idGen,
	}
}
