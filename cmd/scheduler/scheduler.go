package main

import (
	"context"
	"time"

	"booking-schedule/internal/pkg/scheduler"
	"flag"
	"log"

	_ "go.uber.org/automaxprocs"
)

var configType, pathConfig string

func init() {
	flag.StringVar(&configType, "configtype", "file", "type of configuration: environment variables (env) or env/yaml file (file)")
	flag.StringVar(&pathConfig, "config", "./configs/scheduler_config.yml", "path to scheduler config file")
	time.Local = time.UTC
}

func main() {
	flag.Parse()
	ctx := context.Background()
	app, err := scheduler.NewApp(ctx, configType, pathConfig)
	if err != nil {
		log.Fatalf("failed to create scheduler app object:%s\n", err.Error())
	}

	err = app.Run(ctx)
	if err != nil {
		log.Fatalf("failed to run scheduler app: %s", err.Error())
	}
}
