package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Location string `json:"location"`
	Bio      string `json:"bio"`
	Web      string `json:"web"`
	Picture  string `json:"picture"`
	IsRegistered int `json:"is_registered"`
	IsActivated int `json:"is_activated"`
}

type RegisterForm struct {
	FullName string `json:"full_name"`
	Email string `json:"email"`
	Password string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type UserVerificationCode struct {
	gorm.Model
	UserID int `json:"user_id"`
	Code string `json:"code"`
	IsUsed int `json:"is_used"`
}

type Email struct {
	To string `json:"to"`
	From string `json:"from"`
	Subject string `json:"subject"`
	PlainContent string `json:"plain_content"`
	HtmlContent string `json:"html_content"`
}

type UserToken struct {
	gorm.Model

}
