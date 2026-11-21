package handler

import (
	"net/http"
)

func (h *Handler) ID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}
	id := r.PathValue("id")
	url, ok := h.sh.GetURL(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
