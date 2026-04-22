package request

import (
	core_errors "N1ktarchik/Board_of_issues/internal/core/errors"
	"encoding/json"
	"io"
	"net/http"
)

func DecodeAndValidateRequest(r *http.Request, userData any) error {
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		return core_errors.ServerError()
	}

	defer func() { _ = r.Body.Close() }()

	if len(reqData) == 0 {
		return core_errors.BadRequest()
	}

	if err := json.Unmarshal(reqData, &userData); err != nil {
		return core_errors.ServerError()
	}

	return nil
}
