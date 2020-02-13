package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ardiantirta/go-user-management/models"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"net/http"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func Response(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-API-ClientID", "{{APIClientID}}")
	_ = json.NewEncoder(w).Encode(data)
}

func GenerateVerificationCode() string {
	code := uuid.New().String()
	return code
}

func SendVerificationByEmail(user *models.User, verificationCode string) error {
	from := mail.NewEmail("Example User", "test@example.com")
	subject := "Email Verification: user-management-go"
	to := mail.NewEmail(user.FullName, user.Email)
	plainContent := "Please verify your email"
	htmlContent := fmt.Sprintf(`<a href="localhost:3000/auth/verification/%s">verify here</a>`, verificationCode)
	message := mail.NewSingleEmail(from, subject, to, plainContent, htmlContent)
	client := sendgrid.NewSendClient(viper.GetString("sendgrid.api"))
	_, err := client.Send(message)
	if err != nil {
		return errors.New("error when send email verification")
	}

	return nil
}