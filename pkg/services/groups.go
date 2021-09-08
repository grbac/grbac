package services

import (
	"strings"

	"github.com/grbac/grbac/pkg/graph"
)

type MemberError struct {
	member string
	field  string
	err    string
}

func (e *MemberError) Error() string {
	return e.member + ": " + e.field + ": " + e.err
}

func members(members []graph.Member) ([]string, error) {
	var list []string
	for _, member := range members {
		if len(member.Group) != 0 {
			if isGroup(member.Group) {
				list = append(list, toGroupMember(member.Group))
				continue
			}

			return nil, &MemberError{
				member: member.Group,
				field:  "Group",
				err:    "invalid member type",
			}
		}

		if len(member.Subject) != 0 {
			if isAllUsers(member.Subject) {
				list = append(list, "allUsers")
				continue
			}

			if isServiceAccount(member.Subject) {
				list = append(list, toServiceAccountMember(member.Subject))
				continue
			}

			if isUser(member.Subject) {
				list = append(list, toUserMember(member.Subject))
				continue
			}

			return nil, &MemberError{
				member: member.Subject,
				field:  "Subject",
				err:    "invalid member type",
			}
		}

		return nil, &MemberError{
			member: "<nil>",
			field:  "<nil>",
			err:    "member is not set",
		}
	}

	return list, nil
}

func isUserMember(name string) bool {
	return strings.HasPrefix(name, "user:")
}

func isServiceAccountMember(name string) bool {
	return strings.HasPrefix(name, "serviceAccount:")
}

func isGroupMember(name string) bool {
	return strings.HasPrefix(name, "group:")
}

func isAllUsersMember(name string) bool {
	return name == "allUsers"
}

func isGroup(name string) bool {
	return strings.HasPrefix(name, "groups/")
}

func toUserName(name string) string {
	return "users/" + strings.TrimPrefix(name, "user:")
}

func toServiceAccountName(name string) string {
	return "serviceAccounts/" + strings.TrimPrefix(name, "serviceAccount:")
}

func toGroupName(name string) string {
	return "groups/" + strings.TrimPrefix(name, "group:")
}

func toUserMember(name string) string {
	return "user:" + strings.TrimPrefix(name, "users/")
}

func toServiceAccountMember(name string) string {
	return "serviceAccount:" + strings.TrimPrefix(name, "serviceAccounts/")
}

func toGroupMember(name string) string {
	return "group:" + strings.TrimPrefix(name, "groups/")
}
