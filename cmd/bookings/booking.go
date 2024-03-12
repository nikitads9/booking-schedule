package main

import (
	"context"
	"time"

	//_ "go.uber.org/automaxprocs"

	//_ "booking-schedule/cmd/server/docs"
	"booking-schedule/internal/pkg/bookings"
	"flag"
	"log"
)

var pathConfig, pathCert, pathKey string

func init() {
	flag.StringVar(&pathConfig, "config", "./configs/booking_config.yml", "path to config file")
	flag.StringVar(&pathCert, "certfile", "cert.pem", "certificate PEM file")
	flag.StringVar(&pathKey, "keyfile", "key.pem", "key PEM file")
	time.Local = time.UTC
}

//	@title			booking-schedule API
//	@version		1.0
//	@description	This is a service for writing and reading booking entries.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Nikita Denisenok
//	@contact.url	https://vk.com/ndenisenok

//	@license.name	GNU 3.0
//	@license.url	https://www.gnu.org/licenses/gpl-3.0.ru.html

// @host			127.0.0.1:3000
// @BasePath		/bookings
//
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization
//
//	 @Schemes 		http https
//		@Tags			bookings
//
// @tag.name bookings
// @tag.description operations with bookings, suites and intervals
// @tag.name users
// @tag.description operations with user profile such as sign in, sign up, getting profile editing or deleting it
func main() {
	flag.Parse()

	ctx := context.Background()

	app, err := bookings.NewApp(ctx, pathConfig, pathCert, pathKey)
	if err != nil {
		log.Fatalf("failed to create bookings-api app object:%s\n", err.Error())
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("failed to run bookings-api app: %s", err.Error())
	}
}
