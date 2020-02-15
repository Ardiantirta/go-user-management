package http

import (
	"encoding/json"
	"github.com/ardiantirta/go-user-management/helper"
	"github.com/ardiantirta/go-user-management/middleware"
	"github.com/ardiantirta/go-user-management/models"
	"github.com/ardiantirta/go-user-management/services/user"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type UserHandler struct {
	UserService user.Service
}

func NewUserHandler(r *mux.Router, userService user.Service) {
	handler := UserHandler{
		UserService: userService,
	}

	v1 := r.PathPrefix("/me").Subrouter()

	v1.Use(middleware.JwtAuthentication)
	v1.Use(middleware.TwoFactorAuthentication)

	v1.Handle("", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.GetInfo))).Methods(http.MethodGet)
	v1.Handle("", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.UpdateBasicInfo))).Methods(http.MethodPost)
	v1.Handle("/email", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.GetEmailAddress))).Methods(http.MethodGet)
	v1.Handle("/email", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.UpdateEmailAddress))).Methods(http.MethodPost)
	v1.Handle("/password", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.ChangePassword))).Methods(http.MethodPost)
	v1.Handle("/picture", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.SetProfilePicture))).Methods(http.MethodPost)
	v1.Handle("/picture", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.DeleteProfilePicture))).Methods(http.MethodDelete)
	v1.Handle("/tfa", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.TwoFactorAuthenticationStatus))).Methods(http.MethodGet)
	v1.Handle("/tfa/enroll", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.TwoFactorAuthenticationSetup))).Methods(http.MethodGet)
	v1.Handle("/tfa/enroll", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.ActivateTwoFactorAuthentication))).Methods(http.MethodPost)
	v1.Handle("/tfa/remove", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.RemoveTwoFactorAuthentication))).Methods(http.MethodPost)
	v1.Handle("/events", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.ListEventData))).Methods(http.MethodGet)
	v1.Handle("/delete", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.DeleteAccount))).Methods(http.MethodPost)
	v1.Handle("/session", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.SessionLists))).Methods(http.MethodGet)
	v1.Handle("/session", handlers.LoggingHandler(os. Stdout, http.HandlerFunc(handler.DeleteSession))).Methods(http.MethodDelete)
	v1.Handle("/session/other", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.DeleteOtherSessions))).Methods(http.MethodDelete)
	v1.Handle("/session/refresh/token", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.RefreshToken))).Methods(http.MethodGet)
	v1.Handle("/session/access/token", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(handler.NewAccessToken))).Methods(http.MethodGet)
}

func (u *UserHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "something is missing, try re-login"))
		return
	}

	response, err := u.UserService.GetInfo(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (u *UserHandler) UpdateBasicInfo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "something is missing, try re-login"))
		return
	}

	formData := new(models.UpdateUserInfoForm)
	if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "invalid json body"))
		return
	}

	if len(formData.FullName) < 3 || len(formData.FullName) > 128 {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "full_name must be between 3 and 128 chars"))
		return
	}

	if len(formData.Location) > 128 {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "location cannot more than 128 chars"))
		return
	}

	if len(formData.Bio) > 255 {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "bio cannot more than 255 chars"))
		return
	}

	if len(formData.Web) > 128 {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "web cannot more than 128 chars"))
		return
	}

	response, err := u.UserService.UpdateBasicInfo(id, formData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (u *UserHandler) GetEmailAddress(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "something is missing, try re-login"))
		return
	}

	response, err := u.UserService.GetEmailAddress(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (u *UserHandler) UpdateEmailAddress(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "something is missing, try re-login"))
		return
	}

	formData := new(models.ForgotPasswordForm)

	if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "invalid json body"))
		return
	}

	response, err := u.UserService.UpdateEmailAddress(id, formData.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (u *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "something is missing, try re-login"))
		return
	}

	formData := new(models.ChangePasswordForm)
	if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "invalid json body"))
		return
	}

	if len(formData.PasswordCurrent) < 6 || len(formData.PasswordCurrent) > 128 {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "password_current must between 6 and 128 chars"))
		return
	}

	if len(formData.Password) < 6 || len(formData.Password) > 128 {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "password must between 6 and 128 chars"))
		return
	}

	if formData.Password != formData.PasswordConfirm {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "password and password_confirm must be equal"))
		return
	}

	response, err := u.UserService.ChangePassword(id, formData.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (u *UserHandler) SetProfilePicture(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "something is missing, try re-login"))
		return
	}

	link, err := helper.UploadImageToImgur(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "failed to upload image"))
		return
	}

	response, err := u.UserService.SetProfilePicture(id, link)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (u *UserHandler) DeleteProfilePicture(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "something is missing, try re-login"))
		return
	}

	response, err := u.UserService.DeleteProfilePicture(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (u *UserHandler) TwoFactorAuthenticationStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "something is missing, try re-login"))
		return
	}

	response, err := u.UserService.TwoFactorAuthenticationStatus(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (u *UserHandler) TwoFactorAuthenticationSetup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "something is missing, try re-login"))
		return
	}

	response, err := u.UserService.TwoFactorAuthenticationSetup(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (u *UserHandler) ActivateTwoFactorAuthentication(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "something is missing, try re-login"))
		return
	}

	formData := new(models.ActivateTFAForm)
	if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "invalid json body"))
		return
	}

	if len(formData.Secret) != 20 {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "secret must be 20 chars"))
		return
	}

	if len(formData.Code) != 6 {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "code must 6 chars"))
		return
	}

	response, err := u.UserService.ActivateTwoFactorAuthentication(id, formData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (u *UserHandler) RemoveTwoFactorAuthentication(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "something is missing, try re-login"))
		return
	}

	formData := new(models.RemoveTFAForm)
	if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "invalid json body"))
		return
	}

	if len(formData.Password) < 6 || len(formData.Password) > 128 {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "password must between 6 and 128 chars"))
		return
	}

	response, err := u.UserService.RemoveTwoFactorAuthentication(id, formData.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (u *UserHandler) ListEventData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "something is missing, try re-login"))
		return
	}

	response, err := u.UserService.ListEventData(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (u *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "something is missing, try re-login"))
		return
	}

	formData := new(models.RemoveTFAForm)
	if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "invalid json body"))
		return
	}

	if len(formData.Password) < 6 || len(formData.Password) > 128 {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "password must between 6 and 128 chars"))
		return
	}

	response, err := u.UserService.DeleteAccount(id, formData.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (u *UserHandler) SessionLists(w http.ResponseWriter, r *http.Request) {

}

func (u *UserHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "something is missing, try re-login"))
		return
	}

	tokenString := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)

	response, err := u.UserService.DeleteSession(id, tokenString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (u *UserHandler) DeleteOtherSessions(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "something is missing, try re-login"))
		return
	}

	tokenString := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)

	response, err := u.UserService.DeleteOtherSessions(id, tokenString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (u *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "something is missing, try re-login"))
		return
	}

	response, err := u.UserService.RefreshToken(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}

func (u *UserHandler) NewAccessToken(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.Header.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helper.Response(w, helper.ErrorMessage(0, "something is missing, try re-login"))
		return
	}

	response, err := u.UserService.NewAccessToken(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	helper.Response(w, response)
	return
}
