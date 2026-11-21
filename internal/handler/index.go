package handler

import (
	"errors"
	"io"
	"log"
	"net/http"
	"shortener/internal/service"
	"strings"
)

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	bodyString := strings.TrimSpace(string(bodyBytes))
	code, err := h.sh.Shorten(bodyString)
	if err != nil {
		if errors.Is(err, service.ErrInvalidURL) {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}
		log.Printf("error shortening url: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.FormatURL(code)))
}
