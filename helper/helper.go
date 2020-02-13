package helper

import (
	"encoding/json"
	"errors"
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

func SendVerificationByEmail(sgEmail *models.SendGridEmail) error {
	from := sgEmail.From
	subject := sgEmail.Subject
	to := sgEmail.To
	plainContent := sgEmail.PlainContent
	htmlContent := sgEmail.HtmlContent
	message := mail.NewSingleEmail(from, subject, to, plainContent, htmlContent)
	client := sendgrid.NewSendClient(viper.GetString("sendgrid.api"))
	_, err := client.Send(message)
	if err != nil {
		return errors.New("error when send email verification")
	}

	return nil
}