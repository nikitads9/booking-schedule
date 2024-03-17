package main

import (
	"booking-schedule/internal/pkg/sender"
	"context"
	"flag"
	"log"
	"time"

	_ "go.uber.org/automaxprocs"
)

var configType, pathConfig string

func init() {
	flag.StringVar(&configType, "configtype", "file", "type of configuration: environment variables (env) or env/yaml file (file)")
	flag.StringVar(&pathConfig, "config", "./configs/sender_config.yml", "path to sender config file")
	time.Local = time.UTC
}

func main() {
	flag.Parse()
	ctx := context.Background()
	app, err := sender.NewApp(ctx, configType, pathConfig)
	if err != nil {
		log.Fatalf("failed to create sender app object:%s\n", err.Error())
	}

	err = app.Run(ctx)
	if err != nil {
		log.Fatalf("failed to run sender app: %s", err.Error())
	}
}
