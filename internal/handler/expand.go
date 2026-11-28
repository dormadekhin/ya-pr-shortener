package handler

import (
	"net/http"
)

func (h *Handler) Expand(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	url, ok := h.sh.GetURL(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
