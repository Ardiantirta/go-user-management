package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
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
