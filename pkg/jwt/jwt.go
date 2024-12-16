package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/exp/rand"
)

// NewAccessToken returns a new access token
func NewAccessToken(guid, ip, secret string) string {
	token := jwt.New(jwt.SigningMethodHS512)
	token.Claims = jwt.MapClaims{
		"sub": guid,
		"ip":  ip,
	}
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

// NewRefreshToken returns a new refresh token
func NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(uint64(time.Now().Unix()))
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
