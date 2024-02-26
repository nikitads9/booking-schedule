package main

import (
	"context"
	//_ "event-schedule/cmd/server/docs"
	"event-schedule/internal/pkg/events"
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

//	@license.name	GNU 3.0
//	@license.url	https://www.gnu.org/licenses/gpl-3.0.ru.html

//		@host			127.0.0.1:3000
//		@BasePath		/events
//	 @Schemes 		http https
//		@Tags			events
func main() {
	flag.Parse()
	ctx := context.Background()
	app, err := events.NewApp(ctx, pathConfig)
	if err != nil {
		log.Fatalf("failed to create app err:%s\n", err.Error())
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
