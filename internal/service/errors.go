package service

import "errors"

var (
	errLoginAlreadyTaken = errors.New("this login is already taken, try something different :)")
	errInvalidCredentials = errors.New("invalid credentials")
)
