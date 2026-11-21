package handler

import "fmt"

type Shortener interface {
	Shorten(string) (string, error)
	GetURL(string) (string, bool)
}

type Handler struct {
	sh      Shortener
	baseUrl string
}

func NewHandler(shortener Shortener, baseUrl string) *Handler {
	return &Handler{
		sh:      shortener,
		baseUrl: baseUrl,
	}
}

func (h *Handler) FormatURL(code string) string {
	return fmt.Sprintf("%s/%s", h.baseUrl, code)
}
