package setup

import (
	"fmt"
	"github.com/ardiantirta/go-user-management/models"
	"net/url"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func DBConnection() *gorm.DB {
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUser := viper.GetString("database.user")
	dbPass := viper.GetString("database.pass")
	dbName := viper.GetString("database.name")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("sslmode", "disable")
	connStr = fmt.Sprintf("%s?%s", connStr, val.Encode())

	dbConn, err := gorm.Open("postgres", connStr)
	if err != nil {
		logrus.Error(err)
	}

	dbConn.Debug().AutoMigrate(
		&models.User{},
		&models.UserVerificationCode{},
		&models.UserToken{},
		&models.BackUpCode{},
	)

	return dbConn
}
