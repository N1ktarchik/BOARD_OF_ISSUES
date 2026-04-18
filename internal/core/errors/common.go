package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorApp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ErrorApp) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *ErrorApp) StatusCode() int {
	switch e.Code {
	case "TOKEN_NOT_VALID", "JWT_METHOD_NOT_VALID", "UNAUTHORIZED":
		return http.StatusUnauthorized
	case "USER_HAVE_NOT_ACCES", "USER_IS_NOT_OWNER", "USER_IS_NOT_OWNER_OF_TASK":
		return http.StatusForbidden
	case "INVALID_PASSWORD", "PASSWORD_IS_SHORT", "PASSWORD_IS_LONG", "BAD_REQUEST":
		return http.StatusBadRequest
	case "USER_ALREADY_REGISTERED", "EMAIL_ALREADY_REGISTERED":
		return http.StatusConflict
	case "USER_NOT_FOUND", "DESK_NOT_FOUND", "TASK_NOT_FOUND":
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func IsErrorApp(err error) (*ErrorApp, bool) {
	var appErr *ErrorApp
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

func JWTMethodError() *ErrorApp {
	return &ErrorApp{
		Code:    "JWT_METHOD_NOT_VALID",
		Message: "JWT method must be HS256",
	}
}

func JWTTokenNotValid() *ErrorApp {
	return &ErrorApp{
		Code:    "TOKEN_NOT_VALID",
		Message: "JWT token not valid",
	}
}

func UserAlreadyRegistered(login, email string) *ErrorApp {
	return &ErrorApp{
		Code:    "USER_ALREADY_REGISTERED",
		Message: fmt.Sprintf("user with login %s or with email %s already registered", login, email),
	}
}

func InvalidPassword() *ErrorApp {
	return &ErrorApp{
		Code:    "INVALID_PASSWORD",
		Message: "invalid password",
	}
}

func TooShortPassword() *ErrorApp {
	return &ErrorApp{
		Code:    "PASSWORD_IS_SHORT",
		Message: "The password is too short",
	}
}

func TooLongPassword() *ErrorApp {
	return &ErrorApp{
		Code:    "PASSWORD_IS_LONG",
		Message: "The password is too long",
	}
}

func UserNotOwnerOfDesk(userID, deskID string) *ErrorApp {
	return &ErrorApp{
		Code:    "USER_IS_NOT_OWNER",
		Message: fmt.Sprintf("User with id %s is not owner of desk with id %s", userID, deskID),
	}
}

func UserHaveNotAccessToDesk(userID, deskID string) *ErrorApp {
	if len(deskID) != 0 {
		return &ErrorApp{
			Code:    "USER_HAVE_NOT_ACCES",
			Message: fmt.Sprintf("User with id %s have not acces to desk with id %s", userID, deskID),
		}
	}

	return &ErrorApp{
		Code:    "USER_HAVE_NOT_ACCES",
		Message: fmt.Sprintf("User with id %s have not acces to desk", userID),
	}

}

func UserNotOwnerOfTask(userID, taskID string) *ErrorApp {
	return &ErrorApp{
		Code:    "USER_IS_NOT_OWNER_OF_TASK",
		Message: fmt.Sprintf("User with id %s is not owner of taks with id %s", userID, taskID),
	}
}

func EmailAlreadyRegistered(email string) *ErrorApp {
	return &ErrorApp{
		Code:    "EMAIL_ALREADY_REGISTERED",
		Message: fmt.Sprintf("user with email %s already registered", email),
	}
}

func ServerError() *ErrorApp {
	return &ErrorApp{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "internal server error",
	}
}

func BadRequest() *ErrorApp {
	return &ErrorApp{
		Code:    "BAD_REQUEST",
		Message: "bad request",
	}
}

func UnAuthorized() *ErrorApp {
	return &ErrorApp{
		Code:    "UNAUTHORIZED",
		Message: "you have been unaothorized. Please login.",
	}
}

func UserNotFound() *ErrorApp {
	return &ErrorApp{
		Code:    "USER_NOT_FOUND",
		Message: "user not found",
	}
}

func DeskNotFound() *ErrorApp {
	return &ErrorApp{
		Code:    "DESK_NOT_FOUND",
		Message: "desk not found",
	}
}

func TaskNotFound() *ErrorApp {
	return &ErrorApp{
		Code:    "TASK_NOT_FOUND",
		Message: "task not found",
	}
}
