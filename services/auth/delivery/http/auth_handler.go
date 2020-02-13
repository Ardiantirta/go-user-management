package http

import (
	"encoding/json"
	"github.com/ardiantirta/go-user-management/helper"
	"github.com/ardiantirta/go-user-management/middleware"
	"github.com/ardiantirta/go-user-management/models"
	"github.com/ardiantirta/go-user-management/services/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
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
	v1.Handle("/verification/send", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.SendVerificationCode))).Methods(http.MethodPost)
	v1.Handle("/verification/{code}", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.Verification))).Methods(http.MethodPost)
	v1.Handle("/login", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.Login))).Methods(http.MethodPost)
	v1.Handle("/tfa/verify", handlers.LoggingHandler(os.Stdout, middleware.JwtAuthentication(http.HandlerFunc(handler.TwoFactorAuthVerify)))).Methods(http.MethodPost)
	v1.Handle("/tfa/bypass", handlers.LoggingHandler(os.Stdout, middleware.JwtAuthentication(http.HandlerFunc(handler.TwoFactorAuthByPass)))).Methods(http.MethodPost)
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

func (a *AuthHandler) SendVerificationCode(w http.ResponseWriter, r *http.Request) {
	formData := new(models.ResendVerificationForm)

	if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, map[string]interface{}{"code": 0, "message": "invalid json body"})
		return
	}

	params := map[string]interface{}{
		"type": formData.Type,
		"recipient": formData.Recipient,
	}

	response, err := a.AuthService.SendVerificationCode(params)
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
	auth := new(models.AuthenticationForm)

	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, map[string]interface{}{"code": 0, "message": "Invalid json body"})
	}

	response, err := a.AuthService.Login(auth.Email, auth.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (a *AuthHandler) TwoFactorAuthVerify(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "something is missing, please re-login"))
		return
	}

	formData := new(models.VerifyTFAForm)

	if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "invalid json body"))
		return
	}

	if len(formData.Code) != 6 {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "wrong code, try again"))
		return
	}

	response, err := a.AuthService.TwoFactorAuthVerify(id, formData.Code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (a *AuthHandler) TwoFactorAuthByPass(w http.ResponseWriter, r *http.Request) {

}

func (a *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	formData := new(models.ForgotPasswordForm)

	if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, map[string]interface{}{"code": 0, "message": "invalid json body"})
		return
	}

	response, err := a.AuthService.ForgotPassword(formData.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (a *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	formData := new(models.ResetPasswordForm)

	if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, map[string]interface{}{"code": 0, "message": "invalid json body"})
		return
	}

	if len(formData.Password) < 6 || len(formData.Password) > 128 {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, map[string]interface{}{"code": 0, "message": "password length must between 6 and 128 chars"})
		return
	}

	if formData.Password != formData.PasswordConfirm {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, map[string]interface{}{"code": 0, "message": "password and password_confirm must be equal"})
		return
	}

	claims, err := middleware.VerifyToken(formData.Token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, map[string]interface{}{"code": 0, "message": "token is invalid"})
		return
	}

	email := claims.(jwt.MapClaims)["email"].(string)

	response, err := a.AuthService.ResetPassword(email, formData.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}
