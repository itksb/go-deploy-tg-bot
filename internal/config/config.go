package config

import (
	"go-deploy-tg-bot/internal/telegram"
	"log"
	"os"
	"strconv"
	"strings"
)

type (
	Config struct {
		TelegramConfig telegram.Config
	}
)

func NewConfig() (Config, error) {
	return Config{
		TelegramConfig: *telegram.NewConfig(),
	}, nil
}

func (cfg *Config) UseEnv() {
	token, ok := os.LookupEnv("TELEGRAM_TOKEN")
	if ok {
		cfg.TelegramConfig.Token = token
	} else {
		log.Fatal("config error: TELEGRAM_TOKEN env is not setup")
	}

	chatIDs, ok := os.LookupEnv("CHAT_IDS")
	if ok {
		strs := strings.SplitN(chatIDs, ",", 5)
		intChatIDs := make([]int64, len(strs))
		for _, s := range strs {
			d, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				log.Fatal("config error: CHAT_IDS env is not correct. Format: int64,int64")
			}
			intChatIDs = append(intChatIDs, d)
		}
		cfg.TelegramConfig.ChatIDs = intChatIDs
	} else {
		log.Fatal("config error: CHAT_IDS env is not setup")
	}

	cmdDeploy, ok := os.LookupEnv("CMD_DEPLOY")
	if ok {
		cfg.TelegramConfig.CmdToDeploy = cmdDeploy
	} else {
		log.Fatal("config error: CMD_DEPLOY env is not setup")
	}

	debug, ok := os.LookupEnv("DEBUG")
	if ok {
		cfg.TelegramConfig.Debug = bool(strings.ToLower(debug) == "true")
	}

}
