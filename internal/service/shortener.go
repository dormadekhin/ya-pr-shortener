package service

import (
	"errors"
	"shortener/internal/model"
)

var ErrInvalidURL = errors.New("invalid url")
var ErrMaxAttempts = errors.New("max attempts shorten exceeded")

const codeLength = 8
const maxAttempts = 5

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
	for range maxAttempts {
		newCode, err := generateCode(codeLength)
		if err != nil {
			return "", err
		}
		code, ok := s.st.Add(model.Link{Code: newCode, Href: href})
		if ok {
			return code, nil
		}
	}
	return "", ErrMaxAttempts
}

func (s *Shortener) GetURL(code string) (string, bool) {
	url, ok := s.st.Get(code)
	return url, ok
}
