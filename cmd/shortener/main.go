package main

import (
	"fmt"
	"log"
	"net/http"
	"shortener/internal/handler"
	"shortener/internal/repository"
	"shortener/internal/router"
	"shortener/internal/service"
)

func main() {
	port := ":8080"
	baseUrl := fmt.Sprintf("http://localhost%s", port)
	store := repository.NewStore()
	shortener := service.NewShortener(store)
	routerHandler := handler.NewHandler(shortener, baseUrl)
	err := http.ListenAndServe(port, router.NewRouter(routerHandler))
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
