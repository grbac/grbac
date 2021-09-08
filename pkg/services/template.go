package services

import (
	"bufio"
	"bytes"
	"regexp"
	"text/template"
)

var (
	regexAlphaNumeric = regexp.MustCompile("[^A-Za-z0-9]+")

	defaultFuncMap = template.FuncMap{
		"AlphaNumVar": replaceAlphaNumeric,

		"IsUser":           isUserMember,
		"IsServiceAccount": isServiceAccountMember,
		"IsGroup":          isGroupMember,
		"IsAllUsers":       isAllUsersMember,

		"ToUserName":           toUserName,
		"ToServiceAccountName": toServiceAccountName,
		"ToGroupName":          toGroupName,
		"ToPermissionName":     toPermissionName,
	}
)

func replaceAlphaNumeric(name string) string {
	return regexAlphaNumeric.ReplaceAllString(name, "_")
}

func ExecuteTemplate(t *template.Template, data interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	if err := t.Execute(writer, data); err != nil {
		return nil, err
	}

	if err := writer.Flush(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
