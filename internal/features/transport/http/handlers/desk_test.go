package handlers

import (
	dn "Board_of_issuses/internal/core/domains"
	"Board_of_issuses/internal/features/transport/http/dto"
	"Board_of_issuses/internal/features/transport/mocks"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandleCreateDesk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().
		CreateDesk(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, desk *dn.Desk) error {
			if len(desk.Password) <= 6 || len(desk.Password) > 30 {
				return newErrorApp("400", "INVALID_PASSWORD")
			}
			return nil
		}).
		AnyTimes()

	handlers := NewUserHandler(mockService)
	handler := http.HandlerFunc(handlers.HandleCreateDesk)

	requests := []struct {
		password string
		name     string
		id       int
		status   int
	}{
		{"password123", "name1", 123, http.StatusCreated},
		{"pass", "name1", 123, http.StatusInternalServerError},
		{"passworddddddddddddddddddddddddddddddddddddddddd>30.", "name1", 123, http.StatusInternalServerError},
		{"passsword1234", "n", 123, http.StatusBadRequest},
		{"", "n", 123, http.StatusBadRequest},
		{"passsword1234", "", 123, http.StatusBadRequest},
		{"", "", 123, http.StatusBadRequest},
	}

	for _, v := range requests {

		reqBody := &dto.Desk{
			Name:     v.name,
			Password: v.password,
			OwnerId:  v.id,
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("error parse json body")
		}

		req := httptest.NewRequest("POST", "/desks", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(t.Context(), "userID", v.id)
		req = req.WithContext(ctx)

		resp := httptest.NewRecorder()

		handler.ServeHTTP(resp, req)

		if !assert.Equal(t, v.status, resp.Code) {
			msg := resp.Body.String()
			t.Error(msg)
			t.Error("test", v)
			continue
		}

	}
}

func TestHandleChangeDeskName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().
		ChangeDeskName(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, name string, deskId, userID int) error {
			return nil
		}).
		AnyTimes()

	handlers := NewUserHandler(mockService)
	router := mux.NewRouter()
	router.HandleFunc("/desks/{id}/name", handlers.HandleChangeDeskName).Methods("PATCH")

	requests := []struct {
		name   string
		id     int
		deskID int
		code   int
	}{
		{"name1", 123, 123, http.StatusOK},
		{"n", 123, 123, http.StatusBadRequest},
		{"name", 123, 0, http.StatusBadRequest},
		{"name", 123, -123, http.StatusBadRequest},
		{"", 123, 123, http.StatusBadRequest},
	}

	for _, v := range requests {

		reqBody := &dto.UpdateDeskNameRequest{
			Name: v.name,
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("error parce json body")
		}

		url := fmt.Sprintf("/desks/%d/name", v.deskID)
		req := httptest.NewRequest("PATCH", url, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(t.Context(), "userID", v.id)
		req = req.WithContext(ctx)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if !assert.Equal(t, v.code, resp.Code) {
			msg := resp.Body.String()
			t.Error(msg)
			t.Error("test values: ", v)
		}
	}

}

func TestHandleChangeDeskPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().
		ChangeDeskPassword(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, password string, deskID, userID int) error {
			if len(password) <= 6 || len(password) > 30 {
				return newErrorApp("400", "INVALID_PASSWORD")
			}
			return nil
		}).
		AnyTimes()

	handlers := NewUserHandler(mockService)
	router := mux.NewRouter()
	router.HandleFunc("/desks/{id}/password", handlers.HandleChangeDeskPassword).Methods("PATCH")

	requests := []struct {
		password string
		id       int
		deskID   int
		status   int
	}{
		{"password123", 123, 123, http.StatusOK},
		{"pass", 123, 123, http.StatusInternalServerError},
		{"passworddddddddddddddddddddddddddddddddddddddddd>30.", 123, 123, http.StatusInternalServerError},
		{"passwored12", 123, 0, http.StatusBadRequest},
		{"password123", 123, -123, http.StatusBadRequest},
		{"", 123, 0, http.StatusBadRequest},
	}

	for _, v := range requests {

		reqBody := &dto.UpdateDeskPasswordRequest{
			Password: v.password,
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("error parse json body")
		}

		url := fmt.Sprintf("/desks/%d/password", v.deskID)
		req := httptest.NewRequest("PATCH", url, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(t.Context(), "userID", v.id)
		req = req.WithContext(ctx)

		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if !assert.Equal(t, v.status, resp.Code) {
			msg := resp.Body.String()
			t.Error(msg)
			t.Error("test:", v)
			continue
		}

	}
}

func TestHandleChangeDeskOwner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().
		ChangeDeskOwner(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, deskId, userID, newOwnerID int) error {
			return nil
		}).
		AnyTimes()

	handlers := NewUserHandler(mockService)
	router := mux.NewRouter()
	router.HandleFunc("/desks/{id}/owner", handlers.HandleChangeDeskOwner).Methods("PATCH")

	requests := []struct {
		newOwnerID int
		id         int
		deskID     int
		code       int
	}{
		{123, 123, 123, http.StatusOK},
		{0, 123, 123, http.StatusBadRequest},
		{-1, 123, 123, http.StatusBadRequest},
		{123, 123, 0, http.StatusBadRequest},
		{123, 123, -123, http.StatusBadRequest},
	}

	for _, v := range requests {

		reqBody := &dto.UpdateDeskOwnerRequest{
			ID: v.newOwnerID,
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("error parce json body")
		}

		url := fmt.Sprintf("/desks/%d/owner", v.deskID)
		req := httptest.NewRequest("PATCH", url, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(t.Context(), "userID", v.id)
		req = req.WithContext(ctx)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if !assert.Equal(t, v.code, resp.Code) {
			msg := resp.Body.String()
			t.Error(msg)
			t.Error("test values: ", v)
		}
	}

}

func TestHandleDeleteDesk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().
		DeleteDesk(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, deskId, userID int) error {
			return nil
		}).
		AnyTimes()

	handlers := NewUserHandler(mockService)
	router := mux.NewRouter()
	router.HandleFunc("/desks/{id}", handlers.HandleDeleteDesk).Methods("DELETE")

	requests := []struct {
		deskID int
		id     int
		code   int
	}{
		{123, 123, http.StatusNoContent},
		{0, 123, http.StatusBadRequest},
		{-1, 123, http.StatusBadRequest},
	}

	for _, v := range requests {

		url := fmt.Sprintf("/desks/%d", v.deskID)
		req := httptest.NewRequest("DELETE", url, nil)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(t.Context(), "userID", v.id)
		req = req.WithContext(ctx)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if !assert.Equal(t, v.code, resp.Code) {
			msg := resp.Body.String()
			t.Error(msg)
			t.Error("test values: ", v)
		}
	}

}

func TestHandleGetAllDesksId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)

	expectedDesksID := []int{1, 2, 3, 4, 5}

	mockService.EXPECT().
		GetAllDesks(gomock.Any(), gomock.Any()).
		Return(expectedDesksID, nil).
		AnyTimes()

	handlers := NewUserHandler(mockService)
	router := mux.NewRouter()
	router.HandleFunc("/desks", handlers.HandleGetAllDesksId).Methods("GET")

	requests := []struct {
		userID int
		code   int
	}{
		{123, http.StatusOK},
	}

	for _, v := range requests {

		req := httptest.NewRequest("GET", "/desks", nil)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(t.Context(), "userID", v.userID)
		req = req.WithContext(ctx)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if !assert.Equal(t, v.code, resp.Code) {
			msg := resp.Body.String()
			t.Error(msg)
			t.Error("test values: ", v)
		}

		if !assert.Equal(t, v.code, resp.Code) {
			t.Error("test values: ", v)
			return
		}

		if v.code == http.StatusOK {
			var actualDesksID []int
			err := json.Unmarshal(resp.Body.Bytes(), &actualDesksID)
			require.NoError(t, err)
			assert.Equal(t, expectedDesksID, actualDesksID)
		}

	}
}
