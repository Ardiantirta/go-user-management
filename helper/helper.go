package helper

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/ardiantirta/go-user-management/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/viper"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var FormatRFC8601 = "2006-01-02T15:04:05Z"

func Response(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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

const numbers = "1234567890"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func GenerateSecretCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateTFACode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = numbers[seededRand.Intn(len(numbers))]
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

func UploadImageToImgur(r *http.Request) (string, error) {
	url := "https://api.imgur.com/3/image"
	clientID := viper.GetString("imgur.client_key")
	method := "POST"

	file, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	client := &http.Client{}
	req, _ := http.NewRequest(method, url, file)
	req.Header.Add("Authorization", "Client-ID "+clientID)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	m := make(map[string]interface{})
	_ = json.Unmarshal(body, &m)

	return m["data"].(map[string]interface{})["link"].(string), nil
}

func RefreshToken(currentUser *models.User) *models.UserToken {
	userToken := new(models.UserToken)

	isTfa := false
	if currentUser.IsTFA == 1 {
		isTfa = true
	}
	createdAt := time.Now().UTC()
	expiredAt := createdAt.AddDate(0, 0, 7)
	signKey := []byte(viper.GetString("jwt.signkey"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.CustomClaims{
		ID:    int(currentUser.ID),
		Email: currentUser.Email,
		IsTFA: isTfa,
		Code: currentUser.TFACode,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt.Unix(),
		},
	})

	tokenString, _ := token.SignedString(signKey)

	userToken.Token = tokenString
	userToken.Type = "Bearer"
	userToken.DeletedAt = &expiredAt

	return userToken
}