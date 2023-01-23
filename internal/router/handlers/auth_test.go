package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticationHandler_HandleSignUp(t *testing.T) {
	route := "/sign-up"

	t.Run("already-exists", func(t *testing.T) {
		bodyData, _ := json.Marshal(map[string]string{
			"email":    "test-user-1@test.com",
			"password": "password",
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, route, strings.NewReader(string(bodyData)))
		req.Header.Add("Context-type", "application/json")

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "user already exists\n", w.Body.String())
	})

	t.Run("new-user", func(t *testing.T) {
		bodyData, _ := json.Marshal(map[string]string{
			"email":    "test-user-3@test.com",
			"password": "password",
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, route, strings.NewReader(string(bodyData)))
		req.Header.Add("Content-type", "application/json")

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotNil(t, w.Body.String())
		assert.Equal(t, "{\"status\":200,\"body\":{\"message\":\"user created\"}}\n", w.Body.String())
	})
}
