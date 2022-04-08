package telegram

// Config для телеграма из yaml.
type Config struct {
	Token string `mapstructure:"token"`
	Dev   bool   `mapstructure:"dev"`
}
