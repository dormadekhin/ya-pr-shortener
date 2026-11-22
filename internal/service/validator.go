package service

import "net/url"

func isValidURL(str string) bool {
	if len(str) < 1 {
		return false
	}
	u, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}
	return (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}
