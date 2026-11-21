package service

import (
	"errors"
	"regexp"
	"shortener/internal/model"
)

var ErrInvalidURL = errors.New("invalid url")
var codeRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

const codeLength = 8

type Store interface {
	Add(link model.Link) (string, bool)
	Get(code string) (string, bool)
}

type Shortener struct {
	st Store
}

func NewShortener(s Store) *Shortener {
	return &Shortener{st: s}
}

func (s *Shortener) Shorten(href string) (string, error) {
	if !isValidURL(href) {
		return "", ErrInvalidURL
	}
	for range 5 {
		newCode, err := generateCode(codeLength)
		if err != nil {
			return "", err
		}
		code, added := s.st.Add(model.Link{Code: newCode, Href: href})
		if added || len(code) > 0 {
			return code, nil
		}
	}
	return "", errors.New("server error")
}

func (s *Shortener) GetURL(code string) (string, bool) {
	if !codeRegex.MatchString(code) {
		return "", false
	}
	url, ok := s.st.Get(code)
	return url, ok
}
