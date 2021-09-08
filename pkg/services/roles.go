package services

import "strings"

func isRole(name string) bool {
	return strings.HasPrefix(name, "roles/")
}
