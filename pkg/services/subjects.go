package services

import "strings"

func isSubject(name string) bool {
	return isUser(name) || isServiceAccount(name)
}

func isUser(name string) bool {
	return strings.HasPrefix(name, "users/")
}

func isServiceAccount(name string) bool {
	return strings.HasPrefix(name, "serviceAccounts/")
}

const allUsers = "system/allUsers"

func isAllUsers(name string) bool {
	return name == allUsers
}
