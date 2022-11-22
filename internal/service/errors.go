package service

import "errors"

var (
	ErrUserAlreadyExist   = errors.New("user already exist with given mailId")
	ErrInvalidCredentials = errors.New("mail or password are incorrect")
)
