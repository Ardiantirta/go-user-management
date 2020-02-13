package user

import "github.com/ardiantirta/go-user-management/models"

type Service interface {
	GetInfo(id int) (map[string]interface{}, error)
	UpdateBasicInfo(id int, data *models.UpdateUserInfoForm) (map[string]interface{}, error)
	GetEmailAddress(id int) (map[string]interface{}, error)
	UpdateEmailAddress(id int, email string) (map[string]interface{}, error)
	ChangePassword(id int, password string) (map[string]interface{}, error)
	SetProfilePicture(id int) (map[string]interface{}, error)
	DeleteProfilePicture(id int) (map[string]interface{}, error)
	TwoFactorAuthenticationStatus(id int) (map[string]interface{}, error)
	TwoFactorAuthenticationSetup(id int) (map[string]interface{}, error)
	ActivateTwoFactorAuthentication(id int, tfa *models.ActivateTFAForm) (map[string]interface{}, error)
	RemoveTwoFactorAuthentication(id int, password string) (map[string]interface{}, error)
	ListEventData(id int) (map[string]interface{}, error)
	DeleteAccount(id int, password string) (map[string]interface{}, error)
}