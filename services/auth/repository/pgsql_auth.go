package repository

import (
	"errors"
	"github.com/ardiantirta/go-user-management/helper"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"strings"

	"github.com/ardiantirta/go-user-management/models"
	"github.com/ardiantirta/go-user-management/services/auth"
)

type PgsqlAuthRepository struct {
	Conn *gorm.DB
}

func (p *PgsqlAuthRepository) Validate(req *models.RegisterForm) (map[string]interface{}, bool) {
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

func (p *PgsqlAuthRepository) Register(req *models.RegisterForm) (*models.User, *models.UserVerificationCode, error) {
	if resp, ok := p.Validate(req); !ok {
		return nil, nil, errors.New(resp["message"].(string))
	}

	user := new(models.User)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	user.Email = strings.ToLower(req.Email)
	user.FullName = req.FullName
	user.IsRegistered = 0
	user.IsActivated = 0

	if err := p.Conn.Create(&user).Error; err != nil {
		return nil, nil, errors.New("error when create user")
	}

	verificationCode := new(models.UserVerificationCode)
	verificationCode.UserID = int(user.ID)
	verificationCode.Code = helper.GenerateVerificationCode()
	verificationCode.IsUsed = 0

	if err := p.Conn.Create(&verificationCode).Error; err != nil {
		return nil, nil, errors.New("error when create verification code")
	}

	return user, verificationCode, nil
}

func (p *PgsqlAuthRepository) Verification(params map[string]interface{}) error {
	verificationCode := new(models.UserVerificationCode)

	code := params["verification_code"].(string)

	if err := p.Conn.Table("user_verification_codes").
		Where("code = ?", code).
		First(&verificationCode).
		Update("is_used", 1).Error; err != nil {
			return errors.New("verification failed")
	}

	if err := p.Conn.Table("users").
		Where("id = ?", verificationCode.UserID).
		Update("is_registered", 1).Error; err != nil {
			return errors.New("verification failed 2")
	}

	return nil
}

func (p *PgsqlAuthRepository) Login(email, password string) (*models.User, error) {
	panic("implement me")
}

func (p *PgsqlAuthRepository) TwoFactorAuthVerify() error {
	panic("implement me")
}

func (p *PgsqlAuthRepository) TwoFactorAuthByPass() error {
	panic("implement me")
}

func (p *PgsqlAuthRepository) ForgotPassword() error {
	panic("implement me")
}

func (p *PgsqlAuthRepository) ResetPassword() error {
	panic("implement me")
}

func NewPgsqlAuthRepository(conn *gorm.DB) auth.Repository {
	return &PgsqlAuthRepository{
		Conn: conn,
	}
}
