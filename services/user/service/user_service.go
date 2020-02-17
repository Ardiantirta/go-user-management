package service

import (
	"errors"
	"fmt"
	"github.com/ardiantirta/go-user-management/helper"
	"github.com/ardiantirta/go-user-management/models"
	"github.com/ardiantirta/go-user-management/services/user"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService struct {
	UserRepository user.Repository
}

func (u UserService) GetInfo(id int) (map[string]interface{}, error) {
	currentUser, err := u.UserRepository.FetchUserByID(id)
	if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	createdAt := currentUser.CreatedAt.Format(helper.FormatRFC8601)

	return map[string]interface{}{
		"user": map[string]interface{}{
			"id": currentUser.ID,
			"full_name": currentUser.FullName,
			"location": currentUser.Location,
			"bio": currentUser.Bio,
			"web": currentUser.Web,
			"picture": currentUser.Picture,
			"created_at": createdAt,
		},
	}, nil
}

func (u UserService) UpdateBasicInfo(id int, data *models.UpdateUserInfoForm) (map[string]interface{}, error) {
	currentUser, err := u.UserRepository.FetchUserByID(id)
	if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	currentUser.FullName = data.FullName
	currentUser.Location = data.Location
	currentUser.Bio = data.Bio
	currentUser.Web = data.Web

	if err := u.UserRepository.SaveUser(currentUser); err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	return map[string]interface{}{
		"status": true,
	}, nil
}

func (u UserService) GetEmailAddress(id int) (map[string]interface{}, error) {
	currentUser, err := u.UserRepository.FetchUserByID(id)
	if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	return map[string]interface{}{
		"email": currentUser.Email,
	}, nil
}

func (u UserService) UpdateEmailAddress(id int, email string) (map[string]interface{}, error) {
	currentUser, err := u.UserRepository.FetchUserByID(id)
	if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	verificationCode, err := u.UserRepository.CreateVerificationCode(id)
	if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	sg := new(models.SendGridEmail)
	sg.From = mail.NewEmail("User Example 1", "user1@example.com")
	sg.To = mail.NewEmail(currentUser.FullName, email)
	sg.Subject = "Update Email Verification"
	sg.PlainContent = "please verify your email"
	sg.HtmlContent = fmt.Sprintf(`<a href="http://localhost:3000/auth/verification/%s">email verification</a>`, verificationCode.Code)
	if err := helper.SendVerificationByEmail(sg); err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	currentUser.IsActive = 0
	currentUser.IsVerified = 0
	err = u.UserRepository.SaveUser(currentUser); if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	return map[string]interface{}{
		"status": true,
	}, nil
}

func (u UserService) ChangePassword(id int, password string) (map[string]interface{}, error) {
	currentUser, err := u.UserRepository.FetchUserByID(id)
	if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	currentUser.Password = string(hashedPassword)

	if err := u.UserRepository.SaveUser(currentUser); err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	return map[string]interface{}{
		"status": true,
	}, nil
}

func (u UserService) SetProfilePicture(id int, link string) (map[string]interface{}, error) {
	currentUser, err := u.UserRepository.FetchUserByID(id)
	if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	currentUser.Picture = link

	if err := u.UserRepository.SaveUser(currentUser); err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	return map[string]interface{}{
		"status": true,
	}, nil
}

func (u UserService) DeleteProfilePicture(id int) (map[string]interface{}, error) {
	currentUser, err := u.UserRepository.FetchUserByID(id)
	if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	currentUser.Picture = ""

	if err := u.UserRepository.SaveUser(currentUser); err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	return map[string]interface{}{
		"status": true,
	}, nil
}

func (u UserService) TwoFactorAuthenticationStatus(id int) (map[string]interface{}, error) {
	currentUser, err := u.UserRepository.FetchUserByID(id)
	if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	enabled := false
	if currentUser.IsTFA == 1 {
		enabled = true
		enabledAt := currentUser.TFAActivation.Format(helper.FormatRFC8601)
		return map[string]interface{}{
			"enabled": enabled,
			"enabled_at": enabledAt,
		}, nil

	} else {
		enabled = false
		return map[string]interface{}{
			"enabled": enabled,
		}, nil
	}
}

func (u UserService) TwoFactorAuthenticationSetup(id int) (map[string]interface{}, error) {
	currentUser, err := u.UserRepository.FetchUserByID(id)
	if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	secretCode := helper.GenerateSecretCode(20)
	TFACode := helper.GenerateTFACode(6)
	qrCode := helper.GenerateQRCode(TFACode)

	currentUser.SecretCode = secretCode
	currentUser.TFACode = TFACode
	if err := u.UserRepository.SaveUser(currentUser); err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	return map[string]interface{}{
		"secret": secretCode,
		"qr": qrCode,
	}, nil
}

func (u UserService) ActivateTwoFactorAuthentication(id int, tfa *models.ActivateTFAForm) (map[string]interface{}, error) {
	currentUser, err := u.UserRepository.FetchUserByID(id)
	if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	if currentUser.SecretCode != tfa.Secret {
		err := errors.New("wrong secret_code")
		return helper.ErrorMessage(0, err.Error()), err
	}

	if currentUser.TFACode != tfa.Code {
		err := errors.New("wrong code")
		return helper.ErrorMessage(0, err.Error()), err
	}

	now := time.Now().UTC()

	currentUser.IsTFA = 1
	currentUser.TFAActivation = &now
	if err := u.UserRepository.SaveUser(currentUser); err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	backUpCodes := make([]string, 0)
	backUpCodes = append(backUpCodes, helper.GenerateSecretCode(12), helper.GenerateSecretCode(12))
	if err := u.UserRepository.CreateBackUpCode(id, backUpCodes); err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	return map[string]interface{}{
		"backup_codes": backUpCodes,
	}, nil
}

func (u UserService) RemoveTwoFactorAuthentication(id int, password string) (map[string]interface{}, error) {
	currentUser, err := u.UserRepository.FetchUserByID(id)
	if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	currentUser.TFAActivation = nil
	currentUser.IsTFA = 0
	currentUser.SecretCode = ""
	currentUser.TFACode = ""
	if err := u.UserRepository.SaveUser(currentUser); err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	if err := u.UserRepository.DeleteBackUpCodes(id); err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	return map[string]interface{}{
		"status": true,
	}, nil
}

func (u UserService) ListEventData(id int) (map[string]interface{}, error) {
	currentUser, err := u.UserRepository.FetchUserByID(id)
	if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	return map[string]interface{}{
		"status": true,
		"user": currentUser,
	}, nil
}

func (u UserService) DeleteAccount(id int, password string) (map[string]interface{}, error) {
	currentUser, err := u.UserRepository.FetchUserByID(id)
	if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(password)); err != nil {
		return helper.ErrorMessage(0, "password didn't match"), err
	}

	if err := u.UserRepository.DeleteUser(currentUser); err != nil {
		return helper.ErrorMessage(0,  err.Error()), err
	}

	if err := u.UserRepository.DeleteBackUpCodes(id); err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	if err := u.UserRepository.DeleteUserToken(id); err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	return map[string]interface{}{
		"status": true,
	}, nil
}

func (u UserService) SessionLists(id int) (map[string]interface{}, error) {
	panic("implement me")
}

func (u UserService) DeleteSession(id int, currentToken string) (map[string]interface{}, error) {
	_, err := u.UserRepository.FetchUserByID(id)
	if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	if err := u.UserRepository.DeleteUserTokenByToken(id, currentToken); err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	return map[string]interface{}{
		"status": true,
	}, nil
}

func (u UserService) DeleteOtherSessions(id int, currentToken string) (map[string]interface{}, error) {
	_, err := u.UserRepository.FetchUserByID(id)
	if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	if err := u.UserRepository.DeleteUserToken(id); err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	if err := u.UserRepository.CreateNewToken(id, currentToken); err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	return map[string]interface{}{
		"status": true,
	}, nil
}

func (u UserService) RefreshToken(id int) (map[string]interface{}, error) {
	currentUser, err := u.UserRepository.FetchUserByID(id)
	if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	newUserToken := helper.RefreshToken(currentUser)

	expiredAtStr := newUserToken.DeletedAt.Format(helper.FormatRFC8601)

	return map[string]interface{}{
		"access_token": map[string]interface{} {
			"value": newUserToken.Token,
			"type": newUserToken.Type,
			"expired_at": expiredAtStr,
		},
	}, nil
}

func (u UserService) NewAccessToken(id int, oldToken string) (map[string]interface{}, error) {
	_, err := u.UserRepository.FetchUserByID(id)
	if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	if err := u.UserRepository.DeleteUserTokenByToken(id, oldToken); err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	res, err := u.RefreshToken(id)
	if err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	newToken := res["access_token"].(map[string]interface{})["value"].(string)

	if err := u.UserRepository.CreateNewToken(id, newToken); err != nil {
		return helper.ErrorMessage(0, err.Error()), err
	}

	return res, nil
}

func NewUserService(userRepository user.Repository) user.Service {
	return &UserService{
		UserRepository: userRepository,
	}
}
