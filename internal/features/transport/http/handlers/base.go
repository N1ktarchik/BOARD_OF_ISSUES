package handlers

import (
	"log"
	"net/http"
	"os"
)

func getUserIDFromContext(r *http.Request) int {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		log.Panicf("middleware context error")
	}

	return userID
}

func getHomePage() string {
	data, err := os.ReadFile("index.html")
	if err != nil {
		return "<h1>Страница не найдена</h1>"
	}
	return string(data)
}

func (h *UserHandler) HandleBase(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(getHomePage()))
}
