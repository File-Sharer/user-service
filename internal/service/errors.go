package service

import "errors"

var (
	errLoginAlreadyTaken = errors.New("this login is already taken, try something else :)")
	errInvalidCredentials = errors.New("invalid credentials")
)
