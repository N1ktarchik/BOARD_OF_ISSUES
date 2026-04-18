package transport

import (
	"N1ktarchik/Board_of_issues/internal/core/errors"
	"encoding/json"
	"net/http"
)

type JWTResponse struct {
	Token string `json:"token" example:"any JWT-Token"`
}

type ErrorResponse struct {
	Code    int    `json:"code" example:"0"`
	Error   string `json:"error" example:"any string-code of error"`
	Message string `json:"message,omitempty" example:"any error message"`
}

func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code":500, "error":"INTERNAL_SERVER_ERROR"}`))
		return
	}

	w.WriteHeader(code)
	_, _ = w.Write(resp)
}

func RespondWithJWT(w http.ResponseWriter, code int, jwt string) {
	RespondWithJSON(w, code, JWTResponse{Token: jwt})
}

func RespondWithArray(w http.ResponseWriter, code int, key string, data any) {
	RespondWithJSON(w, code, map[string]any{
		key: data,
	})
}

func RespondWithError(w http.ResponseWriter, err error) {
	if appErr, ok := errors.IsErrorApp(err); ok {
		RespondWithJSON(w, appErr.StatusCode(), ErrorResponse{
			Code:    appErr.StatusCode(),
			Error:   appErr.Code,
			Message: appErr.Message,
		})
		return
	}

	RespondWithJSON(w, http.StatusInternalServerError, ErrorResponse{
		Code:    http.StatusInternalServerError,
		Error:   "INTERNAL_SERVER_ERROR",
		Message: "Something went wrong",
	})
}
