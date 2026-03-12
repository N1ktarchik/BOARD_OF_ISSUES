package handlers

import (
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
		{"", "t1@mail.ru", "t1pass", "", http.StatusBadRequest},
		{"t1", "", "t1pass", "", http.StatusBadRequest},
		{"t1", "t1@mail.ru", "", "", http.StatusBadRequest},
		{"", "", "", "", http.StatusBadRequest},
	}

	for _, v := range requests {

		reqBody := &dto.User{
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
		{"", "", "", http.StatusBadRequest},
	}

	for _, v := range requests {

		reqBody := &dto.User{
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

func TestHandleChangeName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().
		ChangeUserName(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, name string, userID int) error {
			return nil
		}).
		AnyTimes()

	handlers := NewUserHandler(mockService)
	handler := http.HandlerFunc(handlers.HandleChangeUserName)

	requests := []struct {
		name string
		id   int
		code int
	}{
		{"name1", 123, http.StatusOK},
		{"n", 123, http.StatusBadRequest},
		{"", 123, http.StatusBadRequest},
	}

	for _, v := range requests {

		reqBody := &dto.UpdateNameRequest{
			Name: v.name,
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("error parce json body")
		}

		req := httptest.NewRequest("PATCH", "/users/name", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(t.Context(), "userID", v.id)
		req = req.WithContext(ctx)

		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		if !assert.Equal(t, v.code, resp.Code) {
			msg := resp.Body.String()
			t.Error(msg)
			t.Error("test values: ", v)
		}
	}

}

func TestHandleChangeEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().
		ChangeUserEmail(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, email string, userID int) error {
			return nil
		}).
		AnyTimes()

	handlers := NewUserHandler(mockService)
	handler := http.HandlerFunc(handlers.HandleChangeUserEmail)

	requests := []struct {
		email string
		id    int
		code  int
	}{
		{"email1@m.ru", 123, http.StatusOK},
		{"email123456@", 123, http.StatusBadRequest},
		{"email123456.", 123, http.StatusBadRequest},
		{"mail@.", 123, http.StatusBadRequest},
		{"email", 123, http.StatusBadRequest},
		{"e", 123, http.StatusBadRequest},
		{"", 123, http.StatusBadRequest},
	}

	for _, v := range requests {

		reqBody := &dto.UpdateEmailRequest{
			Email: v.email,
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("error parce json body")
		}

		req := httptest.NewRequest("PATCH", "/users/email", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(t.Context(), "userID", v.id)
		req = req.WithContext(ctx)

		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		if !assert.Equal(t, v.code, resp.Code) {
			msg := resp.Body.String()
			t.Error(msg)
			t.Error("test values: ", v)
		}
	}

}

func TestHandleChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().
		ChangeUserPassword(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, password string, userID int) error {
			if len(password) <= 6 || len(password) > 30 {
				return newErrorApp("400", "INVALID_PASSWORD")
			}
			return nil
		}).
		AnyTimes()

	handlers := NewUserHandler(mockService)
	handler := http.HandlerFunc(handlers.HandleChangeUserPassword)

	requests := []struct {
		password string
		id       int
		status   int
	}{
		{"password123", 123, http.StatusOK},
		{"pass", 123, http.StatusInternalServerError},
		{"passworddddddddddddddddddddddddddddddddddddddddd>30.", 123, http.StatusInternalServerError},
		{"", 123, http.StatusInternalServerError},
	}

	for _, v := range requests {

		reqBody := &dto.UpdatePasswordRequest{
			Password: v.password,
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("error parse json body")
		}

		req := httptest.NewRequest("PATCH", "/users/password", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(t.Context(), "userID", v.id)
		req = req.WithContext(ctx)

		resp := httptest.NewRecorder()

		handler.ServeHTTP(resp, req)

		if !assert.Equal(t, v.status, resp.Code) {
			msg := resp.Body.String()
			t.Error(msg)
			t.Error("test:", v)
			continue
		}

	}
}

func TestHandleConnectUSerToDesk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().
		ConnectUserToDesk(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, userID, deskID int, password string) error {
			if len(password) <= 6 || len(password) > 30 {
				return newErrorApp("400", "INVALID_PASSWORD")
			}
			return nil
		}).
		AnyTimes()

	handlers := NewUserHandler(mockService)
	handler := http.HandlerFunc(handlers.HandleConnectUserToDesk)

	requests := []struct {
		password string
		id       int
		deskID   int
		status   int
	}{
		{"password123", 123, 12, http.StatusCreated},
		{"pass", 123, 12, http.StatusInternalServerError},
		{"passworddddddddddddddddddddddddddddddddddddddddd>30.", 123, 12, http.StatusInternalServerError},
		{"password123", 123, 0, http.StatusBadRequest},
		{"password123", 123, -112, http.StatusBadRequest},
		{"", 123, -112, http.StatusBadRequest},
	}

	for _, v := range requests {

		reqBody := &dto.ConnectUserToDeskRequest{
			ID:       v.deskID,
			Password: v.password,
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("error parse json body")
		}

		req := httptest.NewRequest("PATCH", "/users", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(t.Context(), "userID", v.id)
		req = req.WithContext(ctx)

		resp := httptest.NewRecorder()

		handler.ServeHTTP(resp, req)

		if !assert.Equal(t, v.status, resp.Code) {
			msg := resp.Body.String()
			t.Error(msg)
			t.Error("test:", v)
			continue
		}

	}
}
