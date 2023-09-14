package test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/go-jose/go-jose.v2"
	"gopkg.in/go-jose/go-jose.v2/jwt"
)

type jwtClaims struct {
	jwt.Claims
	Issuer    string           `json:"iss,omitempty"`
	Subject   string           `json:"sub,omitempty"`
	Audience  jwt.Audience     `json:"aud,omitempty"`
	Expiry    *jwt.NumericDate `json:"exp,omitempty"`
	NotBefore *jwt.NumericDate `json:"nbf,omitempty"`
	IssuedAt  *jwt.NumericDate `json:"iat,omitempty"`
	ID        string           `json:"jti,omitempty"`
	Nickname  string           `json:"https://nickname.com"`
}

type AuthServer struct {
	signer   jose.Signer
	server   *http.Server
	issuer   string
	audience []string
	nickname string
}

func NewAuth0Server() *AuthServer {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	algorithm := jose.RS256

	webKey := jose.JSONWebKey{
		Key:       key,
		KeyID:     "kid",
		Algorithm: string(algorithm),
		Use:       "sig",
	}

	signKey := jose.SigningKey{
		Key:       webKey,
		Algorithm: algorithm,
	}

	signer, err := jose.NewSigner(signKey, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		log.Fatal(err)
	}

	serverAddr := fmt.Sprintf("http://127.0.0.1:%d", port)

	oidcUri := "/.well-known/openid-configuration"
	jwksUri := "/.well-known/jwks.json"

	oidcData := struct {
		JwksUri string `json:"jwks_uri"`
	}{
		JwksUri: serverAddr + jwksUri,
	}

	jwksData := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{webKey.Public()},
	}

	oidcReponse, err := json.Marshal(oidcData)
	if err != nil {
		log.Fatal(err)
	}

	jwksResponse, err := json.Marshal(jwksData)
	if err != nil {
		log.Fatal(err)
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case oidcUri:
			w.Write(oidcReponse)
		case jwksUri:
			w.Write(jwksResponse)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	})

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	go func() {
		server.ListenAndServe()
	}()

	return &AuthServer{
		signer:   signer,
		server:   &server,
		issuer:   serverAddr,
		audience: []string{serverAddr + "/userinfo"},
		nickname: "username",
	}
}

func (s *AuthServer) GetIssuer() string {
	return s.issuer
}

func (s *AuthServer) GetAudience() []string {
	return s.audience
}

func (s *AuthServer) GenerateSubject() string {
	return fmt.Sprintf("auth0|%s", strings.ReplaceAll(uuid.NewString(), "-", "")[:24])
}

func (s *AuthServer) GetNickname() string {
	return s.nickname
}

func (s *AuthServer) GenerateJwt(subject string) (string, error) {
	claims := jwtClaims{
		Issuer:   s.GetIssuer(),
		Audience: s.GetAudience(),
		Subject:  subject,
		Nickname: s.GetNickname(),
	}

	return jwt.Signed(s.signer).Claims(claims).CompactSerialize()
}
