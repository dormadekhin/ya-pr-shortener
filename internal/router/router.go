package router

import (
	"net/http"
)

type Handler interface {
	Create(http.ResponseWriter, *http.Request)
	Expand(http.ResponseWriter, *http.Request)
	Fallback(http.ResponseWriter, *http.Request)
}

func NewRouter(h Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /{$}", h.Create)
	mux.HandleFunc("GET /{id}", h.Expand)
	mux.HandleFunc("/", h.Fallback)
	return mux
}
