package repository

import (
	"errors"
	"fmt"
	"github.com/ardiantirta/go-user-management/helper"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"

	"github.com/ardiantirta/go-user-management/models"
	"github.com/ardiantirta/go-user-management/services/auth"
)

type PgsqlAuthRepository struct {
	Conn *gorm.DB
}

func (p PgsqlAuthRepository) Validate(req *models.RegisterForm) (map[string]interface{}, bool) {
	validate := validator.New()

	if err := validate.Var(req.FullName, "required,min=1,max=128"); err != nil {
		return map[string]interface{}{"code": 0, "message": "full_name is required and must between 1 and 128 chars"}, false
	}

	if err := validate.Var(req.Email, "required,email,max=128"); err != nil {
		return map[string]interface{}{"code": 0, "message": "email must be a valid email and not longer than 128 chars"}, false
	}

	if err := validate.Var(req.Password, "required,min=6,max=128"); err != nil {
		return map[string]interface{}{"code": 0, "message": "password must between 6 and 128 chars"}, false
	}

	if err := validate.VarWithValue(req.Password, req.PasswordConfirm, "eqfield"); err != nil {
		return map[string]interface{}{"code": 0, "message": "password and password_confirm must be equal"}, false
	}

	temp := new(models.User)

	if err := p.Conn.Table("users").
		Where("lower(email) = ?", strings.ToLower(req.Email)).
		First(&temp).Error; err != nil && err != gorm.ErrRecordNotFound {
		return map[string]interface{}{"code": 0, "message": "connection error. please retry"}, false
	}

	if temp.Email != "" {
		return map[string]interface{}{"code": 0, "message": "email already exist"}, false
	}

	return nil, true
}

func (p PgsqlAuthRepository) Register(req *models.RegisterForm) (*models.User, *models.UserVerificationCode, error) {
	if resp, ok := p.Validate(req); !ok {
		return nil, nil, errors.New(resp["message"].(string))
	}

	user := new(models.User)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	user.Email = strings.ToLower(req.Email)
	user.FullName = req.FullName
	user.IsVerified = 0
	user.IsActive = 0
	user.IsTFA = 0
	user.TFAActivation = nil

	if err := p.Conn.Create(&user).Error; err != nil {
		return nil, nil, errors.New("error when create user")
	}

	verificationCode := new(models.UserVerificationCode)
	verificationCode.UserID = int(user.ID)
	verificationCode.Code = helper.GenerateRandomCode()
	verificationCode.IsUsed = 0

	if err := p.Conn.Create(&verificationCode).Error; err != nil {
		return nil, nil, errors.New("error when create verification code")
	}

	return user, verificationCode, nil
}

func (p PgsqlAuthRepository) Verification(params map[string]interface{}) error {
	verificationCode := new(models.UserVerificationCode)

	code := params["verification_code"].(string)

	result := p.Conn.Table("user_verification_codes").
		Where("code = ?", code).
		First(&verificationCode)

	if err := result.Error; err != nil {
		return errors.New("verification code not found")
	}

	if verificationCode.IsUsed == 1 {
		return errors.New("verification code is already used")
	}

	if err := result.Update("is_used", 1).Error; err != nil {
		return errors.New("verification failed")
	}

	if err := p.Conn.Table("users").
		Where("id = ?", verificationCode.UserID).
		Updates(map[string]interface{}{"is_verified": 1, "is_active": 1}).Error; err != nil {
		return errors.New("verification failed")
	}

	return nil
}

func (p PgsqlAuthRepository) SendVerificationCode(email string) (*models.User, *models.UserVerificationCode, error) {
	user := new(models.User)

	if err := p.Conn.Table("users").
		Where("lower(email) = ?", email).
		First(&user).Error; err != nil {
		return nil, nil, errors.New("")
	}

	verificationCode := new(models.UserVerificationCode)
	verificationCode.UserID = int(user.ID)
	verificationCode.Code = helper.GenerateRandomCode()
	verificationCode.IsUsed = 0

	if err := p.Conn.Create(&verificationCode).Error; err != nil {
		return nil, nil, errors.New("error when create verification code")
	}

	return user, verificationCode, nil
}

func (p PgsqlAuthRepository) Login(email, password string) (map[string]interface{}, error) {
	user := new(models.User)

	result := p.Conn.Table("users").
		Where("lower(email) = ? ", email).
		First(&user)
	if err := result.Error; err != nil {
		return nil, errors.New("user not found")
	}

	if err := result.Where("is_active = ?", 1).First(&user).Error; err != nil {
		return nil, errors.New("please verify your email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, errors.New("invalid login credentials, please try again")
	}

	isTfa := false
	if user.IsTFA == 1 {
		isTfa = true
	}

	expiredAt := time.Now().UTC().AddDate(0, 0, 7)
	signKey := []byte(viper.GetString("jwt.signkey"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.CustomClaims{
		ID:    int(user.ID),
		Email: user.Email,
		IsTFA: isTfa,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt.Unix(),
		},
	})

	tokenString, _ := token.SignedString(signKey)

	userToken := new(models.UserToken)
	userToken.Type = "Bearer"
	userToken.Token = tokenString
	userToken.UserID = int(user.ID)
	if err := p.Conn.Table("user_tokens").
		Where("user_id = ? and token = ?", userToken.UserID, userToken.Token).
		FirstOrCreate(&userToken).Error; err != nil {
		return nil, errors.New("create token failed")
	}

	strExpiredAt := expiredAt.Format("2006-01-02T15:04:05Z")

	return map[string]interface{}{
		"require_tfa": isTfa,
		"access_token": map[string]interface{}{
			"value":      userToken.Token,
			"type":       userToken.Type,
			"expired_at": strExpiredAt,
		},
	}, nil
}

func (p PgsqlAuthRepository) ForgotPassword(email string) (*models.User, string, error) {
	user := new(models.User)

	if err := p.Conn.Table("users").
		Where("lower(email) = ?", strings.ToLower(email)).
		Where("is_verified = ?", 1).
		Where("is_active = ?", 1).
		First(&user).Error; err != nil {
		return nil, "", errors.New("user not found")
	}

	isTfa := false
	if user.IsTFA == 1 {
		isTfa = true
	}

	expiredAt := time.Now().UTC().AddDate(0, 0, 7)
	signKey := []byte(viper.GetString("jwt.signkey"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.CustomClaims{
		ID:    int(user.ID),
		Email: user.Email,
		IsTFA: isTfa,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt.Unix(),
		},
	})

	tokenString, _ := token.SignedString(signKey)

	return user, tokenString, nil
}

func (p PgsqlAuthRepository) ResetPassword(email, password string) (map[string]interface{}, error) {
	user := new(models.User)

	fmt.Println(email, password)

	result := p.Conn.Table("users").
		Where("lower(email) = ?", strings.ToLower(email)).
		First(&user)

	if err := result.Error; err != nil {
		return nil, errors.New("user not found")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err := result.Update("password", hashedPassword).Error; err != nil {
		return nil, errors.New("reset password failed")
	}

	return map[string]interface{}{"status": true}, nil
}

func (p PgsqlAuthRepository) FetchUserByID(id int) (*models.User, error) {
	user := new(models.User)

	if err := p.Conn.Table("users").
		Where("id = ?", id).
		First(&user).Error; err != nil {
			return nil, errors.New("user not found")
	}

	return user, nil
}

func (p PgsqlAuthRepository) FetchUserToken(userId int) (*models.UserToken, error) {
	token := new(models.UserToken)

	if err := p.Conn.Table("user_tokens").
		Where("user_id = ?", userId).
		First(&token).Error; err != nil {
			return nil, errors.New("get token failed")
	}

	return token, nil
}

func NewPgsqlAuthRepository(conn *gorm.DB) auth.Repository {
	return &PgsqlAuthRepository{
		Conn: conn,
	}
}
