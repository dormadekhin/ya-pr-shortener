package handler

import "fmt"

type Shortener interface {
	Shorten(string) (string, error)
	GetURL(string) (string, bool)
}

type Handler struct {
	sh      Shortener
	baseURL string
}

func NewHandler(shortener Shortener, baseURL string) *Handler {
	return &Handler{
		sh:      shortener,
		baseURL: baseURL,
	}
}

func (h *Handler) FormatURL(code string) string {
	return fmt.Sprintf("%s/%s", h.baseURL, code)
}
