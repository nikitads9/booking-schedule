package main

import (
	"context"
	//_ "event-schedule/cmd/server/docs"
	app "event-schedule/internal/app"
	"flag"
	"log"
)

var pathConfig string

func init() {
	flag.StringVar(&pathConfig, "config", "config.yml", "path to config file")
}

//	@title			event-schedule API
//	@version		1.0
//	@description	This is a service for writing and reading booking entries.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Nikita Denisenok
//	@contact.url	https://vk.com/ndenisenok

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		127.0.0.1:3000
//	@BasePath	/events/

func main() {
	flag.Parse()
	ctx := context.Background()
	err := app.Start(ctx, pathConfig)
	if err != nil {
		log.Fatalf("failed to start app err:%s\n", err.Error())
	}
}

//TODO: go install github.com/swaggo/swag/cmd/swag@latest и еще chi
