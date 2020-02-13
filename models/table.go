package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email         string    `json:"email" gorm:"type:varchar(255)"`
	Password      string    `json:"password" gorm:"type:varchar(255)"`
	FullName      string    `json:"full_name" gorm:"type:varchar(255)"`
	Location      string    `json:"location" gorm:"type:varchar(255)"`
	Bio           string    `json:"bio" gorm:"type:varchar(255)"`
	Web           string    `json:"web" gorm:"type:varchar(255)"`
	Picture       string    `json:"picture" gorm:"type:varchar(255)"`
	IsVerified    int       `json:"is_verified"`
	IsActive      int       `json:"is_active"`
	IsTFA         int       `json:"is_tfa"`
	TFAActivation *time.Time `json:"tfa_activation"`
	SecretCode string `json:"secret_code" gorm:"type:varchar(255)"`
	QRCode string `json:"qr" gorm:"type:text"`
	TFACode string `json:"tfa_code" gorm:"type:varchar(255)"`
}

type UserVerificationCode struct {
	gorm.Model
	UserID int    `json:"user_id"`
	Code   string `json:"code" gorm:"type:varchar(255)"`
	IsUsed int    `json:"is_used"`
}

type UserToken struct {
	gorm.Model
	UserID int    `json:"user_id"`
	Token  string `json:"token" gorm:"type:text"`
	Type   string `json:"type" gorm:"type:varchar(255)"`
}

type BackUpCode struct {
	gorm.Model
	UserID int `json:"user_id"`
	Code string `json:"code" gorm:"type:varchar(255)"`
}
