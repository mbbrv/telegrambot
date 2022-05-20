package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"telegrambot/internal/mysql"
	telegram "telegrambot/internal/telegram/config"
)

type Service struct {
	User   *mysql.User
	Db     *sqlx.DB
	Bot    *tgbotapi.BotAPI
	Config *telegram.Config
}
