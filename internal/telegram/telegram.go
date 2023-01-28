package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"os"
	"os/exec"
)

type (
	Telegram struct {
		cfg    Config
		logger *zap.Logger
	}
)

func NewTelegram(cfg Config, logger *zap.Logger) (*Telegram, error) {
	tg := &Telegram{
		cfg:    cfg,
		logger: logger,
	}

	return tg, nil
}

func execute(script string, command []string, logger *zap.Logger) (bool, error) {

	cmd := &exec.Cmd{
		Path:         script,
		Args:         command,
		Env:          nil,
		Dir:          "",
		Stdin:        nil,
		Stdout:       os.Stdout,
		Stderr:       os.Stderr,
		ExtraFiles:   nil,
		SysProcAttr:  nil,
		Process:      nil,
		ProcessState: nil,
		Err:          nil,
	}

	logger.Info("Executing command " + cmd.String())

	err := cmd.Start()
	if err != nil {
		return false, err
	}

	err = cmd.Wait()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (tg *Telegram) Run() error {
	bot, err := tgbotapi.NewBotAPI(tg.cfg.Token)
	if err != nil {
		tg.logger.Error("Unable to create telegram api instance")
		return err
	}

	bot.Debug = tg.cfg.Debug

	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := bot.GetUpdatesChan(updateConfig)

	// Let's go through each update that we're getting from Telegram.
	for update := range updates {
		// Telegram can send many types of updates depending on what your Bot
		// is up to. We only want to look at messages for now, so we can
		// discard any other updates.
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		if !sliceContains(tg.cfg.ChatIDs, update.Message.Chat.ID) {
			tg.logger.Warn(fmt.Sprintf("Message from chat which is not allowed! %d \n", update.Message.Chat.ID))
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "У меня есть команды. Главная - update и она обновляет сайт. Не вызывайте команду update слишком часто! "
		case "start":
			msg.Text = "Привет! Я бот для обновления сайта санатория Ключи :)"
		case "status":
			msg.Text = "У меня всё ok :)"
		case "update":
			msg.Text = "Выполняю... Выполню - скажу. Не тревожьте меня!"
			if _, err := bot.Send(msg); err != nil {
				// Note that panics are a bad way to handle errors. Telegram can
				// have service outages or network errors, you should retry sending
				// messages or more gracefully handle failures.
				tg.logger.Error(fmt.Sprintf("Error while sending response: %s \n", err.Error()))
			}

			command := []string{
				tg.cfg.CmdToDeploy,
				"arg1=val1",
			}

			tg.logger.Info("Command to exec: " + tg.cfg.CmdToDeploy)

			_, err := execute(tg.cfg.CmdToDeploy, command, tg.logger)
			if err != nil {
				tg.logger.Error(err.Error())
				msg.Text = "Команда обновления не выполнена. Ошибка, сорян."
			} else {
				msg.Text = "Команда обновления выполнена!"
			}

		default:
			msg.Text = "Что вы мне послали? Я не знаю"
		}

		tg.logger.Info(fmt.Sprintf(
			"Message from Chat.ID: %d with text: %s\n",
			update.Message.Chat.ID,
			update.Message.Text,
		))

		// We'll also say that this message is a reply to the previous message.
		// For any other specifications than Chat ID or Text, you'll need to
		// set fields on the `MessageConfig`.
		msg.ReplyToMessageID = update.Message.MessageID

		// Okay, we're sending our message off! We don't care about the message
		// we just sent, so we'll discard it.
		if _, err := bot.Send(msg); err != nil {
			// Note that panics are a bad way to handle errors. Telegram can
			// have service outages or network errors, you should retry sending
			// messages or more gracefully handle failures.
			tg.logger.Error(fmt.Sprintf("Error while sending response: %s \n", err.Error()))
		}
	}

	return nil
}
