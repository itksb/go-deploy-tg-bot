package main

import (
	"go-deploy-tg-bot/internal/app"
	"go-deploy-tg-bot/internal/config"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	cfg.UseEnv()

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer application.Close()

	log.Fatal(application.Run())
}
