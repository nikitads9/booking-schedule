package main

import (
	"context"
	"flag"
	"log"
	"os"

	"event-schedule/internal/pkg/scheduler"
)

var pathConfig string

func init() {
	flag.StringVar(&pathConfig, "config", "config.yml", "path to config file")
}

func main() {
	flag.Parse()
	ctx := context.Background()
	go func() {
		file, err := os.Open(pathConfig)
		if err != nil {
			return
		}
		file.Close()
	}()
	app, err := scheduler.NewApp(ctx, pathConfig)
	if err != nil {
		log.Fatalf("failed to create app err:%s\n", err.Error())
	}

	err = app.Run(ctx)
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
