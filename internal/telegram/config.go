package telegram

type Config struct {
	Token       string
	ChatIDs     []int64
	Debug       bool
	CmdToDeploy string
}

func NewConfig() *Config {
	return &Config{
		Token:   "",
		ChatIDs: make([]int64, 1),
		Debug:   false,
	}
}
