package main

import (
	"context"
	app "event-schedule/internal/app"
	"flag"
	"log"
)

var pathConfig string

func init() {
	flag.StringVar(&pathConfig, "config", "config.yml", "path to config file")
}

func main() {
	flag.Parse()
	ctx := context.Background()
	err := app.Start(ctx, pathConfig)
	if err != nil {
		log.Fatalf("failed to start app err:%s\n", err.Error())
	}
}
