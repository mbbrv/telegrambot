package router

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"log"
	"telegrambot/internal/models"
	"telegrambot/internal/telegram/keyboards"
)

type Router interface {
	Route() (string, error)
}

type router struct {
	user   *models.User
	update *tgbotapi.Update
	bot    *tgbotapi.BotAPI
	db     *sqlx.DB
	auth   bool
}

func NewRouter(user models.User, update tgbotapi.Update, bot tgbotapi.BotAPI, db sqlx.DB, auth bool) Router {
	return router{&user, &update, &bot, &db, auth}
}

func (r router) Route() (string, error) {
	if r.auth {
		if errMsg, err := r.handleGreetings(); err != nil {
			log.Println(err)
			return errMsg, err
		}
	}

	switch r.update.CallbackData() {
	case "daily":
		if errMsg, err := r.handleDaily(); err != nil {
			return errMsg, err
		}
	case "care":
		if errMsg, err := r.handleCare(); err != nil {
			return errMsg, err
		}
	case "appointment":
		if errMsg, err := r.handleAppointment(); err != nil {
			return errMsg, err
		}
	case "description":
		if errMsg, err := r.handleDescriptionEdit(); err != nil {
			return errMsg, err
		}
	case "morningTime":
		if errMsg, err := r.handleTime("morning"); err != nil {
			return errMsg, err
		}
	case "morningSetTime":
		if errMsg, err := r.handleSetTime("morning"); err != nil {
			return errMsg, err
		}
	case "eveningTime":
		if errMsg, err := r.handleTime("evening"); err != nil {
			return errMsg, err
		}
	case "eveningSetTime":
		if errMsg, err := r.handleSetTime("evening"); err != nil {
			return errMsg, err
		}
	default:
		data := keyboards.Data{}
		if err := json.Unmarshal([]byte(r.update.CallbackData()), &data); err == nil {
			errMsg, err := r.handleChangeTime(data)
			if err != nil {
				return errMsg, err
			}
		}
	}

	if r.update.Message.Command() == "description" {
		if errMsg, err := r.handleDescription(); err != nil {
			log.Println(err)
			return errMsg, err
		}
	}

	return "", nil
}
