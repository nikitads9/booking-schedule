package main

import (
	"context"
	"time"

	"event-schedule/internal/pkg/sender"
	"flag"
	"log"
)

var pathConfig string

func init() {
	flag.StringVar(&pathConfig, "config", "./configs/sender_config.yml", "path to sender config file")
	time.Local = time.UTC
}

func main() {
	flag.Parse()
	ctx := context.Background()
	app, err := sender.NewApp(ctx, pathConfig)
	if err != nil {
		log.Fatalf("failed to create sender app object:%s\n", err.Error())
	}

	err = app.Run(ctx)
	if err != nil {
		log.Fatalf("failed to run sender app: %s", err.Error())
	}
}
