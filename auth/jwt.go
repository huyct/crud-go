package jwt

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func Create(username string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 12)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}

func Verify(r *http.Request) error {
	tokenString := Extract(r)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func Extract(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")

	if token != "" {
		return token
	}

	bearerToken := r.Header.Get("Authorization")
	return strings.Split(bearerToken, " ")[1]
}

func EctractUsernameFromToken(r *http.Request) (string, error) {
	var username string
	tokenString := Extract(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err != nil {
		return username, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username = fmt.Sprintf("%v", claims["username"])
	}
	return username, nil
}
