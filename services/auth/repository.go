package auth

import "github.com/ardiantirta/go-user-management/models"

type Repository interface {
	Register(req *models.RegisterForm) (*models.User, *models.UserVerificationCode, error)
	Verification(params map[string]interface{}) error
	Login(email, password string) (*models.User, error)
	TwoFactorAuthVerify() error
	TwoFactorAuthByPass() error
	ForgotPassword() error
	ResetPassword() error
}
