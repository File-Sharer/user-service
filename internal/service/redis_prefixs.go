package service

import "fmt"

const (
	USER_PREFIX = "user:%s" // <userID>
	USER_LOGIN_PREFIX = "user-login:%s" // <login>
)

func UserPrefix(userID string) string {
	return fmt.Sprintf(USER_PREFIX, userID)
}

func UserLoginPrefix(login string) string {
	return fmt.Sprintf(USER_LOGIN_PREFIX, login)
}
