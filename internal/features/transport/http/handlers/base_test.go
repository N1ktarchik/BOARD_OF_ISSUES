package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleBase(t *testing.T) {
	var handlers = NewUserHandler(nil)
	handler := http.HandlerFunc(handlers.HandleBase)

	requests := []struct {
		request string
		status  int
		method  string
	}{
		{"/", http.StatusOK, "GET"},
		{"/", http.StatusOK, "POST"},
		{"/", http.StatusOK, "DELETE"},
		{"/", http.StatusOK, "PUT"},
		{"/", http.StatusOK, "PATCH"},
		{"/", http.StatusOK, "HEAD"},
	}

	for _, v := range requests {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest(v.method, v.request, nil)

		handler.ServeHTTP(resp, req)

		assert.Equal(t, v.status, resp.Code)
	}
}
