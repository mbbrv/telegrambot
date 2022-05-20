// Package cmd /*
package cmd

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
	"log"
	"telegrambot/internal/helpers"
	"telegrambot/internal/mysql"
	"telegrambot/internal/telegram/keyboards"
	router2 "telegrambot/internal/telegram/router"
	"telegrambot/internal/vars"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Запуск бота для взаимодействия с пользователем",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		updateConfig := tgbotapi.NewUpdate(0)

		updateConfig.Timeout = 30

		updates := Service.Bot.GetUpdatesChan(updateConfig)

		for update := range updates {
			message := helpers.GetMessage(&update)
			if errMsg, err := processServe(update); err != nil {
				errorMsg(errMsg, message.Chat.ID, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func processServe(update tgbotapi.Update) (string, error) {
	var auth = false
	message := helpers.GetMessage(&update)
	if message == nil {
		return "", nil
	}
	if message.Text == "/start" {

		msg := tgbotapi.NewMessage(message.Chat.ID, vars.WelcomeMessage)
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{Keyboard: keyboards.GetKeyboardButtonsStart()}
		msg.ParseMode = "HTML"

		if _, err := Bot.Send(msg); err != nil {
			log.Println(err)
			return vars.ErrorDefault, err
		}

		return "", nil
	}

	if !message.IsCommand() && update.CallbackQuery == nil && message.Contact == nil && message.Text != vars.KeyboardButtonUsername {
		del := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
		if _, err := Bot.Send(del); err != nil {
			log.Println(err)

			return "", nil
		}

		return "", nil
	}

	if message.Contact != nil {
		err := mysql.UserEnrichmentByPhoneNumb(Db, message)
		if err == nil {
			log.Println(err)
		}

		auth = true
	}

	if message.Text == vars.KeyboardButtonUsername {
		err := mysql.UserEnrichmentByUsername(Db, message)
		if err == nil {
			log.Println(err)
		}

		auth = true
	}

	//TODO: объединить вместе селекты к юзеру
	if user, ok, err := mysql.IsAuth(Db, message.Chat); ok {
		if update.CallbackQuery == nil {
			if err := user.UpdateFirstName(Db, message); err != nil {
				return "", err
			}
		}

		router := router2.NewRouter(&user, &update, Bot, Db, auth)
		if errMsg, err := router.Route(); err != nil {
			log.Println(err)
			return errMsg, err
		}

	} else {
		log.Println(err)
		return vars.ErrorNoUser, err
	}
	return "", nil
}

func errorMsg(message string, chatId int64, err error) {
	msg := tgbotapi.NewMessage(chatId, message)
	if ConfigTelegram.Dev {
		msg = tgbotapi.NewMessage(chatId, err.Error())
	}

	if _, err := Bot.Send(msg); err != nil {
		log.Println(err)
	}
}
