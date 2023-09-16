package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sesaquecruz/go-chat-broadcaster/test"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJwtMiddleware(t *testing.T) {
	auth := test.NewAuth0Server()

	sub := auth.GenerateSubject()
	jwt, err := auth.GenerateJwt(sub)
	assert.Nil(t, err)

	r := gin.New()
	r.Use(JwtMiddleware(auth.GetIssuer(), auth.GetAudience()))

	r.GET("/", func(c *gin.Context) {
		claims, err := JwtGetClaims(c)
		assert.Nil(t, err)

		assert.Equal(t, auth.GetIssuer(), claims.Issuer)
		assert.Equal(t, auth.GetAudience()[0], claims.Audience[0])
		assert.Equal(t, auth.GetNickname(), claims.Nickname)
		assert.Equal(t, sub, claims.Subject)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+jwt)

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
