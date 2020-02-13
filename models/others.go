package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type CustomClaims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	IsTFA bool `json:"is_tfa"`
	jwt.StandardClaims
}

type SendGridEmail struct {
	From         *mail.Email `json:"from"`
	To           *mail.Email `json:"to"`
	Subject      string      `json:"subject"`
	PlainContent string      `json:"plain_content"`
	HtmlContent  string      `json:"html_content"`
}
