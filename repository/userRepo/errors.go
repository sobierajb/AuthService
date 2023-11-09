package userRepo

import "errors"

var (
	ErrHashPassword  = errors.New("cannot hash password")
	ErrEmptyParam    = errors.New("parameter is 0 string or not defined")
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("user password not matched")
)
