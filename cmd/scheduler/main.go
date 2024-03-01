package main

import (
	"context"
	"time"

	"event-schedule/internal/pkg/scheduler"
	"flag"
	"log"
)

var pathConfig string

func init() {
	flag.StringVar(&pathConfig, "config", "config.yml", "path to config file")
	time.Local = time.UTC
}

func main() {
	flag.Parse()
	ctx := context.Background()
	app, err := scheduler.NewApp(ctx, pathConfig)
	if err != nil {
		log.Fatalf("failed to create app err:%s\n", err.Error())
	}

	err = app.Run(ctx)
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
