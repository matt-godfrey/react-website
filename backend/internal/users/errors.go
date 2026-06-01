package users

import "errors"

var ErrUserAlreadyExists = errors.New("user already exists")

var ErrEmailNotFound = errors.New("email not found")
