package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	auth "github.com/auth0/go-jwt-middleware/v2"
	middleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

type JwtClaims struct {
	Issuer    string
	Subject   string
	Audience  []string
	Expiry    int64
	NotBefore int64
	IssuedAt  int64
	ID        string
	Nickname  string
}

type JwtCustomClaims struct {
	Nickname string `json:"https://nickname.com"`
}

func (c *JwtCustomClaims) Validate(ctx context.Context) error {
	return nil
}

func JwtGetClaims(c *gin.Context) (*JwtClaims, error) {
	claims, ok := c.Request.Context().Value(middleware.ContextKey{}).(*validator.ValidatedClaims)
	if !ok {
		return nil, errors.New("fail to get jwt claims")
	}

	registered := claims.RegisteredClaims
	custom := claims.CustomClaims.(*JwtCustomClaims)

	return &JwtClaims{
		Issuer:    registered.Issuer,
		Subject:   registered.Subject,
		Audience:  registered.Audience,
		Expiry:    registered.Expiry,
		NotBefore: registered.NotBefore,
		IssuedAt:  registered.IssuedAt,
		ID:        registered.ID,
		Nickname:  custom.Nickname,
	}, nil
}

func JwtMiddleware(issuer string, audience []string) gin.HandlerFunc {
	url, err := url.Parse(issuer)
	if err != nil {
		log.Fatal(err)
	}

	provider := jwks.NewCachingProvider(url, 10*time.Minute)

	validator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		url.String(),
		audience,
		validator.WithCustomClaims(func() validator.CustomClaims {
			return &JwtCustomClaims{}
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	middleware := auth.New(
		validator.ValidateToken,
		auth.WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {}),
	)

	return func(c *gin.Context) {
		unauthorized := true

		var authorize http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			unauthorized = false
			c.Request = r
			c.Next()
		}

		middleware.CheckJWT(authorize).ServeHTTP(c.Writer, c.Request)

		if unauthorized {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
