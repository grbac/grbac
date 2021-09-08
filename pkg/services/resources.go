package services

import "net/url"

func isFullResourceName(name string) bool {
	if name == "@animeshon" {
		return true
	}

	if len(name) == 0 || name[:2] != "//" {
		return false
	}

	_, err := url.Parse("https:" + name)
	return err == nil
}
