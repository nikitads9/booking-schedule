package main

import (
	"context"
	"time"

	_ "go.uber.org/automaxprocs"
	//_ "booking-schedule/cmd/auth/docs"

	"booking-schedule/internal/pkg/auth"
	"flag"
	"log"
)

var configType, pathConfig, pathCert, pathKey string

func init() {
	flag.StringVar(&configType, "configtype", "file", "type of configuration: environment variables (env) or env/yaml file (file)")
	flag.StringVar(&pathConfig, "config", "./configs/auth_config.yml", "path to config file")
	flag.StringVar(&pathCert, "certfile", "cert.pem", "certificate PEM file")
	flag.StringVar(&pathKey, "keyfile", "key.pem", "key PEM file")
	time.Local = time.UTC
}

//	@title			auth API
//	@version		1.0
//	@description	This is a basic auth service for booking API.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Nikita Denisenok
//	@contact.url	https://vk.com/ndenisenok

//	@license.name	GNU 3.0
//	@license.url	https://www.gnu.org/licenses/gpl-3.0.ru.html

// @host			127.0.0.1:5000
// @BasePath		/auth

// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization
//
//	 @Schemes 		http https
//		@Tags			auth
//
// @tag.name auth
// @tag.description sign in and sign up operations
func main() {
	flag.Parse()

	ctx := context.Background()

	app, err := auth.NewApp(ctx, configType, pathConfig, pathCert, pathKey)
	if err != nil {
		log.Fatalf("failed to create auth-api app object:%s\n", err.Error())
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("failed to run auth-api app: %s", err.Error())
	}
}
