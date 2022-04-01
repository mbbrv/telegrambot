package telegram

type Config struct {
	Token string `mapstructure:"token"`
	Dev   bool   `mapstructure:"dev"`
}
