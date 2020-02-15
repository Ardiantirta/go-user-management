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
			response = helper.ErrorMessage(0, "missing auth token")
			w.WriteHeader(http.StatusUnauthorized)
			helper.Response(w, response)
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := VerifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			helper.Response(w, helper.ErrorMessage(0, "invalid auth token"))
			return
		}

		id := claims.(jwt.MapClaims)["id"].(float64)
		email := claims.(jwt.MapClaims)["email"].(string)
		isTFA := claims.(jwt.MapClaims)["is_tfa"].(bool)
		code := claims.(jwt.MapClaims)["code"].(string)

		r.Header.Set("id", strconv.Itoa(int(id)))
		r.Header.Set("email", email)
		r.Header.Set("is_tfa", strconv.FormatBool(isTFA))
		r.Header.Set("code", code)

		next.ServeHTTP(w, r)
	})
}

func CheckClientID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientKey := viper.GetString("client.key")
		apiClientID := r.Header.Get("X-API-ClientID")

		if len(apiClientID) < 1 {
			w.WriteHeader(http.StatusForbidden)
			helper.Response(w, helper.ErrorMessage(0, "please provide X-API-ClientID"))
			return
		}

		if apiClientID != clientKey {
			w.WriteHeader(http.StatusBadRequest)
			helper.Response(w, helper.ErrorMessage(0, "wrong X-API-ClientID"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func TwoFactorAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isTfa, err := strconv.ParseBool(r.Header.Get("is_tfa"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			helper.Response(w, helper.ErrorMessage(0, "something is missing, try re-login"))
			return
		}

		code := r.Header.Get("code")

		canAccess := false
		if isTfa == true && len(code) == 6 {
			canAccess = true
		} else if isTfa == false {
			canAccess = true
		}

		if !canAccess {
			w.WriteHeader(http.StatusBadRequest)
			helper.Response(w, helper.ErrorMessage(0, "Please Verify Two Factor Authentication"))
			return
		}

		next.ServeHTTP(w, r)
	})
}