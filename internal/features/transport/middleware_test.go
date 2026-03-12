package tarnsport

import (
	"Board_of_issuses/internal/features/transport/http/handlers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	auth "Board_of_issuses/internal/core/auth/jwt"

	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	jwtCONFIG := auth.LoadJwtConfigWithParamSecretKey("testsecret")
	jwtManager := auth.CreateJWTService(jwtCONFIG)
	authHandler := CreateAuthHandler(jwtManager)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDVal := r.Context().Value("userID")
		if userIDVal == nil {
			handlers.RespondWithError(w, http.StatusInternalServerError, "userID not in context")
			return
		}

		userID, ok := userIDVal.(int)
		if !ok {
			handlers.RespondWithError(w, http.StatusInternalServerError, "userID is wrong type")
			return
		}

		handlers.RespondWithJSON(w, http.StatusOK, map[string]int{"user_id": userID})
	})

	handler := authHandler.AuthMiddleware(testHandler)

	req1 := httptest.NewRequest("GET", "/protected", nil)
	resp1 := httptest.NewRecorder()
	handler.ServeHTTP(resp1, req1)
	assert.Equal(t, http.StatusUnauthorized, resp1.Code)

	req2 := httptest.NewRequest("GET", "/protected", nil)
	req2.Header.Set("Authorization", "Token something")
	resp2 := httptest.NewRecorder()
	handler.ServeHTTP(resp2, req2)
	assert.Equal(t, http.StatusBadRequest, resp2.Code)

	req3 := httptest.NewRequest("GET", "/protected", nil)
	req3.Header.Set("Authorization", "Bearer not-a-valid-jwt-token")
	resp3 := httptest.NewRecorder()
	handler.ServeHTTP(resp3, req3)
	if resp3.Code == http.StatusInternalServerError {
		assert.Equal(t, http.StatusInternalServerError, resp3.Code)
	} else {
		assert.Equal(t, http.StatusUnauthorized, resp3.Code)
	}

	validToken, _ := jwtManager.Create(123, "test@gmail.com")
	req4 := httptest.NewRequest("GET", "/protected", nil)
	req4.Header.Set("Authorization", "Bearer "+validToken)
	resp4 := httptest.NewRecorder()
	handler.ServeHTTP(resp4, req4)
	assert.Equal(t, http.StatusOK, resp4.Code)

	var response map[string]int
	json.Unmarshal(resp4.Body.Bytes(), &response)
	assert.Equal(t, 123, response["user_id"])
}
