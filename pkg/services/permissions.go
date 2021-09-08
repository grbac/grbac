package services

import "strings"

func isPermission(name string) bool {
	return strings.HasPrefix(name, "permissions/")
}

func toPermissionId(name string) string {
	return strings.TrimPrefix(name, "permissions/")
}

func toPermissionName(name string) string {
	return "permissions/" + name
}

// isValidPermissionId enforces the Google Cloud IAM permission format
// [service].[resource].[verb].
func isValidPermissionId(name string) bool {
	return len(strings.Split(toPermissionId(name), ".")) == 3
}
