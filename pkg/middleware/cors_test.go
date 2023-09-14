package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCorsMiddleware(t *testing.T) {
	r := gin.New()
	r.Use(CorsMiddleware())

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/", nil)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "*", res.Header.Get("Access-Control-Allow-Origin"))
}
