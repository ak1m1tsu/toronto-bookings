package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anthdm/weavebox"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticationHandler_HandleSignUp(t *testing.T) {
	store := NewMockMongoUserStore()
	handler := NewAuthenticationHandler(store)
	route := "/sign-up"

	app := weavebox.New()
	app.Post(route, handler.HandleSignUp)

	t.Run("already-exists", func(t *testing.T) {
		bodyData, _ := json.Marshal(map[string]string{
			"email":    "test-user-1@test.com",
			"password": "password",
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, route, strings.NewReader(string(bodyData)))
		req.Header.Add("Context-type", "application/json")

		app.ServeHTTP(w, req)

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

		app.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotNil(t, w.Body.String())
		assert.Equal(t, "{\"status\":200,\"body\":{\"message\":\"user created\"}}\n", w.Body.String())
	})
}
