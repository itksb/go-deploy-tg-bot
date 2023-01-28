package app

import (
	"go-deploy-tg-bot/internal/config"
	"go-deploy-tg-bot/internal/telegram"
	"go.uber.org/zap"
	"io"
	"log"
)

type App struct {
	telegram *telegram.Telegram
	logger   *zap.Logger

	io.Closer
}

func NewApp(cfg config.Config) (*App, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Println("Couldn't create zap logger instance")
		return nil, err
	}
	defer logger.Sync() // flushes buffer, if any

	tgram, err := telegram.NewTelegram(cfg.TelegramConfig, logger)
	if err != nil {
		logger.Error("Unable to create telegram instance")
		return nil, err
	}

	app := App{
		logger:   logger,
		telegram: tgram,
	}

	return &app, nil
}

func (app *App) Run() error {
	err := app.telegram.Run()
	return err
}

func (app *App) Close() error {
	app.logger.Info("closing app")
	return nil
}
