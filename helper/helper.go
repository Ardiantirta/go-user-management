package helper

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/ardiantirta/go-user-management/models"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/viper"
	"math/rand"
	"net/http"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func Response(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-API-ClientID", "{{APIClientID}}")
	_ = json.NewEncoder(w).Encode(data)
}

func ErrorMessage(code int, message string) map[string]interface{} {
	return map[string]interface{}{
		"code": code,
		"message": message,
	}
}

func GenerateRandomCode() string {
	code := uuid.New().String()
	return code
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func GenerateSecretCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateQRCode(content string) string {
	var png []byte
	png, _ = qrcode.Encode(content, qrcode.Medium, 256)
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
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