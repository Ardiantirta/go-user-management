package repository

import (
	"errors"
	"fmt"
	"github.com/ardiantirta/go-user-management/helper"
	"github.com/ardiantirta/go-user-management/models"
	"github.com/ardiantirta/go-user-management/services/user"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

type UserRepository struct {
	Conn *gorm.DB
}

func (u UserRepository) FetchUserByID(id int) (*models.User, error) {
	currentUser := new(models.User)

	if err := u.Conn.Table("users").
		Where("id = ?", id).
		First(&currentUser).Error; err != nil {
			return nil, errors.New("user not found")
	}

	return currentUser, nil
}

func (u UserRepository) SaveUser(user *models.User) error {
	if err := u.Conn.Table("users").
		Where("id = ?", user.ID).
		Save(&user).Error; err != nil {
			return errors.New("failed to update user data")
	}

	return nil
}

func (u UserRepository) CreateVerificationCode(id int) (*models.UserVerificationCode, error) {
	verificationCode := new(models.UserVerificationCode)
	verificationCode.UserID = id
	verificationCode.Code = helper.GenerateRandomCode()
	verificationCode.IsUsed = 0

	if err := u.Conn.Create(&verificationCode).Error; err != nil {
		return nil, errors.New("error when create verification code")
	}

	return verificationCode, nil
}

func (u UserRepository) CreateBackUpCode(id int, codes []string) error {
	query := `insert into back_up_codes (user_id, code, created_at, updated_at) values `
	countInsert := 0
	createdAt := time.Now().Format(helper.FormatRFC8601)
	for _, c := range codes {
		values := fmt.Sprintf("" +
			"(%d, '%s', '%s', '%s'),",
			id, c, createdAt, createdAt,
			)
		query = query + values
		countInsert++
	}

	if countInsert > 0 {
		query = strings.TrimSuffix(query, ",")
		if err := u.Conn.Exec(query).Error; err != nil {
			return errors.New("failed to create backup_codes")
		}
	}

	return nil
}

func (u UserRepository) DeleteBackUpCodes(id int) error {
	if err := u.Conn.Unscoped().Table("back_up_codes").
		Where("user_id = ?", id).
		Delete(models.BackUpCode{}).Error; err != nil {
			return errors.New("failed to delete backup_codes")
	}

	return nil
}

func (u UserRepository) DeleteUser(user *models.User) error {
	if err := u.Conn.Unscoped().Delete(&user).Error; err != nil {
		return errors.New("failed to delete user")
	}

	return nil
}

func (u UserRepository) DeleteUserToken(userID int) error {
	if err := u.Conn.Unscoped().Table("user_tokens").
		Where("user_id = ?", userID).
		Delete(models.UserToken{}).Error; err != nil {
			return errors.New("failed to delete user token")
	}

	return nil
}

func (u UserRepository) DeleteUserTokenByToken(userID int, token string) error {
	if err := u.Conn.Unscoped().Table("user_tokens").
		Where("user_id = ?", userID).
		Where("token = ?", token).
		Delete(models.UserToken{}).Error; err != nil {
			return errors.New("failed to delete user token")
	}

	fmt.Println(userID, token)

	return nil
}

func (u UserRepository) CreateNewToken(id int, newToken string) error {
	token := new(models.UserToken)

	token.Type = "Bearer"
	token.UserID = id
	token.Token = newToken

	if err := u.Conn.Create(&token).Error; err != nil {
		return errors.New("failed to create new token")
	}

	return nil
}

func (u UserRepository) CheckToken(id int, token string) error {
	userToken := new(models.UserToken)

	if err := u.Conn.Table("user_tokens").
		Where("user_id = ?", id).
		Where("token = ?", token).
		First(&userToken).Error; err != nil {
			return errors.New("user token not found")
	}

	return nil
}

func NewUserRepository(conn *gorm.DB) user.Repository {
	return &UserRepository{
		Conn: conn,
	}
}
