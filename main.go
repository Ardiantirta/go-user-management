package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ardiantirta/go-user-management/helper"
	"github.com/ardiantirta/go-user-management/setup"

	authHttp "github.com/ardiantirta/go-user-management/services/auth/delivery/http"
	_authRepository "github.com/ardiantirta/go-user-management/services/auth/repository"
	_authService "github.com/ardiantirta/go-user-management/services/auth/service"

	userHttp "github.com/ardiantirta/go-user-management/services/user/delivery/http"
	_userRepository "github.com/ardiantirta/go-user-management/services/user/repository"
	_userService "github.com/ardiantirta/go-user-management/services/user/service"
)

func init() {
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if viper.GetBool("debug") {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbConn := setup.DBConnection()
	defer func() {
		err := dbConn.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	r := mux.NewRouter()

	r.Handle("/", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"message": "welcome to go-user-management api",
			"code": http.StatusOK,
		}
		helper.Response(w, response)
		return
	}))).Methods(http.MethodGet)

	authRepository := _authRepository.NewPgsqlAuthRepository(dbConn)
	authService := _authService.NewAuthService(authRepository)
	authHttp.NewAuthHandler(r, authService)

	userRepository := _userRepository.NewUserRepository(dbConn)
	userService := _userService.NewUserService(userRepository)
	userHttp.NewUserHandler(r, userService)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "X-API-ClientID"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE"})

	logrus.Fatal(http.ListenAndServe(viper.GetString("server.address"), handlers.CORS(headersOk, originsOk, methodsOk)(r)))
}
