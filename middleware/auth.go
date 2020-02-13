package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"

	"github.com/ardiantirta/go-user-management/helper"
)

func VerifyToken(tokenString string) (jwt.Claims, error) {
	signKey := []byte(viper.GetString("jwt.signkey"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return signKey, err
	})

	if err != nil {
		return nil, err
	}

	return token.Claims, nil
}

func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := make(map[string]interface{})
		tokenString := r.Header.Get("Authorization")

		if len(tokenString) == 0 {
			response = map[string]interface{}{"code": 0, "message": "missing auth token"}
			w.WriteHeader(http.StatusUnauthorized)
			helper.Response(w, response)
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := VerifyToken(tokenString)
		if err != nil {
			response = map[string]interface{}{"code": 0, "message": "invalid auth token"}
			w.WriteHeader(http.StatusForbidden)
			helper.Response(w, response)
			return
		}

		id := claims.(jwt.MapClaims)["id"].(float64)
		email := claims.(jwt.MapClaims)["email"].(string)
		isTFA := claims.(jwt.MapClaims)["is_tfa"].(bool)

		r.Header.Set("id", strconv.Itoa(int(id)))
		r.Header.Set("email", email)
		r.Header.Set("is_tfa", strconv.FormatBool(isTFA))

		next.ServeHTTP(w, r)
	})
}
