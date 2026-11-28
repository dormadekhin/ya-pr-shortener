package handler

import "net/http"

func (h *Handler) Fallback(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Wrong request.", http.StatusBadRequest)
}
