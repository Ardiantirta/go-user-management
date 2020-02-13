package http

import (
	"encoding/json"
	"github.com/ardiantirta/go-user-management/helper"
	"github.com/ardiantirta/go-user-management/models"
	"github.com/ardiantirta/go-user-management/services/auth"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

type AuthHandler struct {
	AuthService auth.Service
}

func NewAuthHandler(r *mux.Router, authService auth.Service) {
	handler := &AuthHandler{
		AuthService: authService,
	}

	v1 := r.PathPrefix("/auth").Subrouter()

	v1.Handle("/register", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.Register))).Methods(http.MethodPost)
	v1.Handle("/verification/{code}", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.Verification))).Methods(http.MethodPost)
	v1.Handle("/login", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.Login))).Methods(http.MethodPost)
	v1.Handle("/tfa/verify", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.TwoFactorAuthVerify))).Methods(http.MethodPost)
	v1.Handle("/tfa/bypass", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.TwoFactorAuthByPass))).Methods(http.MethodPost)
	v1.Handle("/password/forgot", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.ForgotPassword))).Methods(http.MethodPost)
	v1.Handle("/password/reset", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.ResetPassword))).Methods(http.MethodPost)
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	req := new(models.RegisterForm)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, map[string]interface{}{"code": 0, "message": "Invalid json body"})
		return
	}

	response, err := a.AuthService.Register(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (a *AuthHandler) Verification(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]

	params := map[string]interface{}{
		"verification_code": code,
	}

	response, err := a.AuthService.Verification(params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

}

func (a *AuthHandler) TwoFactorAuthVerify(w http.ResponseWriter, r *http.Request) {

}

func (a *AuthHandler) TwoFactorAuthByPass(w http.ResponseWriter, r *http.Request) {

}

func (a *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {

}

func (a *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {

}
