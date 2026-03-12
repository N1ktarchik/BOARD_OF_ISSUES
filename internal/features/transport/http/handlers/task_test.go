package handlers

import (
	dn "Board_of_issuses/internal/core/domains"
	"Board_of_issuses/internal/features/transport/http/dto"
	"Board_of_issuses/internal/features/transport/mocks"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().
		CreateTask(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, task *dn.Task) error {
			return nil
		}).
		AnyTimes()

	handlers := NewUserHandler(mockService)
	handler := http.HandlerFunc(handlers.HandleCreateTask)

	requests := []struct {
		UserID      int
		DeskId      int
		Name        string
		Description string
		Time        time.Time
		status      int
	}{
		{123, 123, "name", "desc", time.Now(), http.StatusCreated},
		{123, 0, "name", "desc", time.Now(), http.StatusBadRequest},
		{123, -1, "name", "desc", time.Now(), http.StatusBadRequest},
		{123, 123, "n", "desc", time.Now(), http.StatusBadRequest},
		{123, 123, "", "desc", time.Now(), http.StatusBadRequest},
	}

	for _, v := range requests {

		reqBody := &dto.Task{
			DeskId:      v.DeskId,
			Name:        v.Name,
			Description: v.Description,
			Time:        v.Time,
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("error parse json body")
		}

		req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(t.Context(), "userID", v.UserID)
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

func TestHandleDeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().
		DeleteTask(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, taskID, userID int) error {
			if userID != 123 {
				return errors.New("user is not owner")
			}
			return nil
		}).
		AnyTimes()

	handlers := NewUserHandler(mockService)
	router := mux.NewRouter()
	router.HandleFunc("/tasks/{id}", handlers.HandleDeleteTask).Methods("DELETE")

	requests := []struct {
		deskID int
		userID int
		status int
	}{
		{123, 123, http.StatusNoContent},
		{0, 123, http.StatusBadRequest},
		{-1, 123, http.StatusBadRequest},
		{1, 1234, http.StatusInternalServerError},
	}

	for _, v := range requests {
		url := fmt.Sprintf("/tasks/%d", v.deskID)
		req := httptest.NewRequest("DELETE", url, nil)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(t.Context(), "userID", v.userID)
		req = req.WithContext(ctx)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if !assert.Equal(t, v.status, resp.Code) {
			msg := resp.Body.String()
			t.Error(msg)
			t.Error("test values: ", v)
		}
	}
}

func TestHandleComplyteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().
		ComplyteTask(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, taskID, userID int) error {
			return nil
		}).
		AnyTimes()

	handlers := NewUserHandler(mockService)
	router := mux.NewRouter()
	router.HandleFunc("/tasks/{id}/complyte", handlers.HandleComplyteTask).Methods("PATCH")

	requests := []struct {
		taskID int
		userID int
		status int
	}{
		{123, 123, http.StatusOK},
		{0, 123, http.StatusBadRequest},
		{-1, 123, http.StatusBadRequest},
	}

	for _, v := range requests {
		url := fmt.Sprintf("/tasks/%d/complyte", v.taskID)
		req := httptest.NewRequest("PATCH", url, nil)
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(t.Context(), "userID", v.userID)
		req = req.WithContext(ctx)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if !assert.Equal(t, v.status, resp.Code) {
			msg := resp.Body.String()
			t.Error(msg)
			t.Error("test values: ", v)
		}
	}
}

func TestHandleAddTimeToTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().
		UpdateTaskTime(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, taskID, userID int, userTime int) error {
			if userID != 123 {
				return errors.New("user is not owner")
			}
			return nil
		}).
		AnyTimes()

	handlers := NewUserHandler(mockService)
	router := mux.NewRouter()
	router.HandleFunc("/tasks/{id}/time", handlers.HandleAddTimeToTask).Methods("PATCH")

	requests := []struct {
		taskID   int
		userID   int
		userTime int
		status   int
	}{
		{123, 123, 1, http.StatusOK},
		{0, 123, 1, http.StatusBadRequest},
		{-1, 123, 1, http.StatusBadRequest},
		{1, 1234, 1, http.StatusInternalServerError},
	}

	for _, v := range requests {
		reqBody := &dto.UpdateTaskTimeRequest{
			Hours: v.userTime,
		}
		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("error parce json body")
		}

		url := fmt.Sprintf("/tasks/%d/time", v.taskID)
		req := httptest.NewRequest("PATCH", url, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(t.Context(), "userID", v.userID)
		req = req.WithContext(ctx)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if !assert.Equal(t, v.status, resp.Code) {
			msg := resp.Body.String()
			t.Error(msg)
			t.Error("test values: ", v)
		}
	}
}

func TestHandleChangeTaskDescription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().
		ChangeTaskDescription(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, taskID, userID int, descrip string) error {
			if userID != 123 {
				return errors.New("user is not owner")
			}
			return nil
		}).
		AnyTimes()

	handlers := NewUserHandler(mockService)
	router := mux.NewRouter()
	router.HandleFunc("/tasks/{id}/description", handlers.HandleChangeTaskDescription).Methods("PATCH")

	requests := []struct {
		taskID int
		userID int
		desc   string
		status int
	}{
		{123, 123, "1", http.StatusOK},
		{123, 123, "", http.StatusOK},
		{0, 123, "22", http.StatusBadRequest},
		{-1, 123, "1", http.StatusBadRequest},
		{1, 1234, "1", http.StatusInternalServerError},
	}

	for _, v := range requests {
		reqBody := &dto.UpdateTaskDescriptionRequest{
			Description: v.desc,
		}
		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			t.Errorf("error parce json body")
		}

		url := fmt.Sprintf("/tasks/%d/description", v.taskID)
		req := httptest.NewRequest("PATCH", url, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(t.Context(), "userID", v.userID)
		req = req.WithContext(ctx)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if !assert.Equal(t, v.status, resp.Code) {
			msg := resp.Body.String()
			t.Error(msg)
			t.Error("test values: ", v)
		}
	}
}

func TestGetAllTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedTasks := []dn.Task{
		{Id: 12, UserId: 123, DeskId: 123, Name: "test1", Description: "testtest1", Done: false},
		{Id: 1234, UserId: 123, DeskId: 1234, Name: "test2", Description: "testtest2", Done: true},
	}

	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().
		GetAllTasks(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, userID, deskID int) ([]dn.Task, error) {
			if userID != 123 {
				return nil, errors.New("user is not owner")
			}
			return expectedTasks, nil
		}).
		AnyTimes()

	handlers := NewUserHandler(mockService)
	router := mux.NewRouter()
	router.HandleFunc("/tasks/{deskId}", handlers.HandleGetAllTasks).Methods("GET")

	requests := []struct {
		userID int
		deskID int
		code   int
	}{
		{123, 123, http.StatusOK},
		{123, 0, http.StatusBadRequest},
		{123, -1, http.StatusBadRequest},
		{12, 123, http.StatusInternalServerError},
	}

	for _, v := range requests {

		url := fmt.Sprintf("/tasks/%d", v.deskID)
		req := httptest.NewRequest("GET", url, nil)
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
			var actualTasks []dn.Task
			err := json.Unmarshal(resp.Body.Bytes(), &actualTasks)
			require.NoError(t, err)
			assert.Equal(t, expectedTasks, actualTasks)
		}

	}

}
