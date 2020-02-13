package service

import (
	"github.com/ardiantirta/go-user-management/helper"
	"github.com/ardiantirta/go-user-management/models"
	"github.com/ardiantirta/go-user-management/services/auth"
)

type AuthService struct {
	AuthRepository auth.Repository
}

func (a *AuthService) Register(req *models.RegisterForm) (map[string]interface{}, error) {
	user, verificationCode, err := a.AuthRepository.Register(req)
	if err != nil {
		return map[string]interface{}{"code": 0, "message": err.Error()}, err
	}

	if err := helper.SendVerificationByEmail(user, verificationCode.Code); err != nil {
		return map[string]interface{}{"code": 0, "message": err.Error()}, err
	}

	mapResponse := map[string]interface{}{"status": true}
	return mapResponse, nil
}

func (a *AuthService) Verification(params map[string]interface{}) (map[string]interface{}, error) {
	err := a.AuthRepository.Verification(params)
	if err != nil {
		return map[string]interface{}{"code": 0, "message": err.Error()}, err
	}

	mapResponse := map[string]interface{}{"status": true}
	return mapResponse, nil
}

func (a *AuthService) Login(email, password string) (map[string]interface{}, error) {
	panic("implement me")
}

func (a *AuthService) TwoFactorAuthVerify() (map[string]interface{}, error) {
	panic("implement me")
}

func (a *AuthService) TwoFactorAuthByPass() (map[string]interface{}, error) {
	panic("implement me")
}

func (a *AuthService) ForgotPassword() (map[string]interface{}, error) {
	panic("implement me")
}

func (a *AuthService) ResetPassword() (map[string]interface{}, error) {
	panic("implement me")
}

func NewAuthService(authRepository auth.Repository) auth.Service {
	return &AuthService{
		AuthRepository: authRepository,
	}
}
