package user

import "github.com/ardiantirta/go-user-management/models"

type Repository interface {
	FetchUserByID(id int) (*models.User, error)
	SaveUser(user *models.User) error
	DeleteUser(user *models.User) error
	DeleteUserToken(id int) error
	DeleteBackUpCodes(id int) error
	CreateVerificationCode(id int) (*models.UserVerificationCode, error)
	CreateBackUpCode(id int, codes []string) error
	DeleteBackUpCode(id int) error
}
