package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type User struct {
	gorm.Model
	Email       string `json:"email"`
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	Location    string `json:"location"`
	Bio         string `json:"bio"`
	Web         string `json:"web"`
	Picture     string `json:"picture"`
	IsVerified  int    `json:"is_verified"`
	IsActive int `json:"is_active"`
}

type RegisterForm struct {
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type ForgotPasswordForm struct {
	Email string `json:"email"`
}

type ResetPasswordForm struct {
	Token string `json:"token"`
	Password string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type ResendVerificationForm struct {
	Type string `json:"type"`
	Recipient string `json:"recipient"`
}

type UserVerificationCode struct {
	gorm.Model
	UserID int    `json:"user_id"`
	Code   string `json:"code"`
	IsUsed int    `json:"is_used"`
}

type Authentication struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserToken struct {
	gorm.Model
	UserID int `json:"user_id"`
	Token string `json:"token" gorm:"type:text"`
	Type string `json:"type"`
}

type CustomClaims struct {
	ID int `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

type SendGridEmail struct {
	From *mail.Email `json:"from"`
	To *mail.Email `json:"to"`
	Subject string `json:"subject"`
	PlainContent string `json:"plain_content"`
	HtmlContent string `json:"html_content"`
}
