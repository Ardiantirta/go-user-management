package auth

import "github.com/ardiantirta/go-user-management/models"

type Repository interface {
	Register(req *models.RegisterForm) (*models.User, *models.UserVerificationCode, error)
	Verification(params map[string]interface{}) error
	SendVerificationCode(email string) (*models.User, *models.UserVerificationCode, error)
	Login(email, password string) (map[string]interface{}, error)
	TwoFactorAuthVerify() (map[string]interface{}, error)
	TwoFactorAuthByPass() (map[string]interface{}, error)
	ForgotPassword(email string) (*models.User, string,  error)
	ResetPassword(email, password string) (map[string]interface{}, error)
}
