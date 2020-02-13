package auth

import "github.com/ardiantirta/go-user-management/models"

type Service interface {
	Register(req *models.RegisterForm) (map[string]interface{}, error)
	Verification(params map[string]interface{}) (map[string]interface{}, error)
	Login(email, password string) (map[string]interface{}, error)
	TwoFactorAuthVerify() (map[string]interface{}, error)
	TwoFactorAuthByPass() (map[string]interface{}, error)
	ForgotPassword() (map[string]interface{}, error)
	ResetPassword() (map[string]interface{}, error)
}