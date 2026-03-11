package handlers

import (
	er "Board_of_issuses/internal/core"
	"Board_of_issuses/internal/core/domains"
	"Board_of_issuses/internal/features/transport/http/dto"
	"Board_of_issuses/internal/features/transport/mocks"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func newErrorApp(code, message string) *er.ErrorApp {
	return &er.ErrorApp{
		Code:    code,
		Message: message,
	}
}

func TestHandleCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().
		Registration(gomock.Any(), gomock.Any()).
		Return("", nil).
		AnyTimes()

	handlers := NewUserHandler(mockService)
	handler := http.HandlerFunc(handlers.HandleCreateUser)

	requests := []struct {
		name     string
		email    string
		password string
		login    string
		status   int
	}{
		{"t1", "t1@mail.ru", "t1pass", "t1log", http.StatusCreated},
		{"", "t1@mail.ru", "t1pass", "t1log", http.StatusBadRequest},
		{"t1", "", "t1pass", "t1log", http.StatusBadRequest},
		{"t1", "t1@mail.ru", "", "t1log", http.StatusBadRequest},
		{"t1", "t1@mail.ru", "t1pass", "", http.StatusBadRequest},
	}

	for _, v := range requests {

		reqBody := dto.User{
			Name:     v.name,
			Email:    v.email,
			Password: v.password,
			Login:    v.login,
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("error parse json body")
		}

		req := httptest.NewRequest("POST", "/register", bytes.NewReader(jsonBody))

		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()

		handler.ServeHTTP(resp, req)

		if !assert.Equal(t, v.status, resp.Code) {
			msg := resp.Body.String()
			t.Error(msg)
			t.Error("test (email,password,login,code):", v)
			continue
		}

		if resp.Code != http.StatusOK {
			continue
		}

		var response map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &response)

		assert.Equal(t, "jwt-token-123", response["access_token"], "токен не совпадает")
		assert.Equal(t, "Bearer", response["token_type"], "неправильный token_type")
	}
}

func TestHandleLoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().
		Authorization(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, user *domains.User) (string, error) {
			// Если пароль пустой - возвращаем ошибку INVALID_PASSWORD
			if len(user.Password) <= 6 || len(user.Password) > 30 {
				return "", newErrorApp("400", "INVALID_PASSWORD")
			}
			return "jwt-token-123", nil
		}).
		AnyTimes()

	handlers := NewUserHandler(mockService)
	handler := http.HandlerFunc(handlers.HandleLoginUser)

	requests := []struct {
		email    string
		password string
		login    string
		status   int
	}{
		{"t1@mail.ru", "t1passs", "t1log", http.StatusOK},
		{"t1@mail.ru", "t1passs", "", http.StatusOK},
		{"", "t1passs", "t1log", http.StatusOK},
		{"t1@mail.ru", "", "t1log", http.StatusBadRequest},
		{"", "t1pass", "", http.StatusBadRequest},
		{"", "t1passsssssssssssssssssssssss>30", "", http.StatusBadRequest},
	}

	for _, v := range requests {

		reqBody := dto.User{
			Email:    v.email,
			Password: v.password,
			Login:    v.login,
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("error parse json body")
		}

		req := httptest.NewRequest("POST", "/register", bytes.NewReader(jsonBody))

		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()

		handler.ServeHTTP(resp, req)

		if !assert.Equal(t, v.status, resp.Code) {
			msg := resp.Body.String()
			t.Error(msg)
			t.Error("test (email,password,login,code):", v)
			continue
		}

		if resp.Code != http.StatusOK {
			continue
		}

		var response map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &response)

		assert.Equal(t, "jwt-token-123", response["access_token"], "токен не совпадает")
		assert.Equal(t, "Bearer", response["token_type"], "неправильный token_type")

	}
}
