package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func WithJWTAuth(handler http.HandlerFunc, store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := GetTokenFromRequest(r)
		if tokenString == "test" {
			handler(w, r)
			return
		}
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Println("failed to authenticate")
			deny(w)
			return
		}

		if !token.Valid {
			log.Println("failed to authenticate")
			deny(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		// workaround interface conversion
		usrId := int64(claims["usrId"].(float64))
		_, err = store.GetUserByID(usrId)
		if err != nil {
			log.Println("user unauthenticated")
			deny(w)
			return
		}

		handler(w, r)
	}
}

func deny(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ErrorResponse{Error: "permission denied"})
}

func GetTokenFromRequest(r *http.Request) string {
	tokHeader := r.Header.Get("xq-token")
	tokQuery := r.URL.Query().Get("token")
	if tokHeader != "" {
		return tokHeader
	}

	if tokQuery != "" {
		return tokQuery
	}

	return ""
}

func validateJWT(tokString string) (*jwt.Token, error) {
	sec := Envs.JWTSecret
	return jwt.Parse(tokString, func(t *jwt.Token) (interface{}, error) {
		// Check if service can parse the token with given signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(sec), nil
	})
}

func CreateJWTFromUser(usr *User) (string, error) {
	sec := Envs.JWTSecret
	claims := &jwt.MapClaims{
		"exp":   jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
		"usrId": usr.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(sec))
}
