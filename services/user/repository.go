package user

import "github.com/ardiantirta/go-user-management/models"

type Repository interface {
	FetchUserByID(id int) (*models.User, error)
	SaveUser(user *models.User) error
	DeleteUser(user *models.User) error
	DeleteUserToken(id int) error
	DeleteUserTokenByToken(id int, token string) error
	DeleteBackUpCodes(id int) error
	CreateVerificationCode(id int) (*models.UserVerificationCode, error)
	CreateBackUpCode(id int, codes []string) error
	CreateNewToken(id int, newToken string) error
}
