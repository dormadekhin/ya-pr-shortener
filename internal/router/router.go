package router

import (
	"net/http"
)

type Handler interface {
	Index(http.ResponseWriter, *http.Request)
	Id(http.ResponseWriter, *http.Request)
}

func NewRouter(h Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Index)
	mux.HandleFunc("/{id}", h.Id)
	return mux
}
