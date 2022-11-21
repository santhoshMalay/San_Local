package service

import "errors"

var (
	ErrUserAlreadyExist = errors.New("User already exist with given mailId")
)
