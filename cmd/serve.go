// Package cmd /*
package cmd

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
	"log"
	"telegrambot/internal/helpers"
	"telegrambot/internal/mysql"
	"telegrambot/internal/vars"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{

	Use:   "serve",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		var auth = false

		Bot.Debug = ConfigTelegram.Dev

		updateConfig := tgbotapi.NewUpdate(0)

		updateConfig.Timeout = 30

		updates := Bot.GetUpdatesChan(updateConfig)

		for update := range updates {

			message, err := helpers.GetMessage(update)
			if err != nil {
				//errorMsg(vars.HandleDefault, message.Chat.ID, bot, err)
				log.Println(err)

				continue
			}

			if message.Text == "/start" {

				msg := tgbotapi.NewMessage(message.Chat.ID, vars.WelcomeMessage)
				msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{Keyboard: helpers.GetKeyboardButtonsStart()}

				if _, err := Bot.Send(msg); err != nil {
					errorMsg(vars.HandleDefault, update.Message.Chat.ID, err)
					log.Println(err)

					continue
				}

				continue
			}

			if !message.IsCommand() && update.CallbackQuery == nil && message.Contact == nil && message.Text != vars.KeyboardButtonUsername {

				//Возможно, ничего не надо отправлять при удалении сообщения
				//msg := tgbotapi.NewMessage(message.Chat.ID, vars.HandleKeyboard)
				//if _, err := bot.Send(msg); err != nil {
				//	errorMsg(message.Chat.ID, bot, err)
				//	log.Println(err)
				//
				//	continue
				//}

				del := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
				if _, err := Bot.Send(del); err != nil {
					log.Println(err)

					continue
				}

				continue
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
				if auth {
					greetingsMsg(message.Chat.ID)
					auth = false
				}

				if update.CallbackData() == "care" {
					err := user.ChangeCareStatus(Db)
					if err != nil {
						errorMsg(vars.HandleDefault, message.Chat.ID, err)
						log.Println(err)

						continue
					}

					textMessage := user.GetChangeCareStatus(vars.CareDisabled, vars.CareEnabled)
					msg := tgbotapi.NewMessage(message.Chat.ID, textMessage)
					if _, err := Bot.Send(msg); err != nil {
						errorMsg(vars.HandleDefault, message.Chat.ID, err)
						log.Println(err)

						continue
					}
				}

				if update.CallbackData() == "appointment" {
					textMessage, err := user.GetPreparedAppointment(Db)
					if err != nil {
						errorMsg(vars.HandleDefault, message.Chat.ID, err)
						log.Println(err)

						continue
					}

					msg := tgbotapi.NewMessage(message.Chat.ID, textMessage)
					msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: helpers.GetInlineButtonsMain()}
					if _, err := Bot.Send(msg); err != nil {
						errorMsg(vars.HandleDefault, message.Chat.ID, err)
						log.Println(err)

						continue
					}
				}

				if message.Command() == "description" {
					descriptionMsg(message.Chat.ID)
				}

			} else {
				errorMsg(vars.HandleNoUser, message.Chat.ID, err)
				log.Println(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
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

func greetingsMsg(chatId int64) {
	msg := tgbotapi.NewMessage(chatId, vars.AuthSuccessMessage)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	if _, err := Bot.Send(msg); err != nil {
		errorMsg(vars.HandleDefault, chatId, err)
		log.Println(err)

		return
	}

	descriptionMsg(chatId)
}

func descriptionMsg(chatId int64) {
	msg := tgbotapi.NewMessage(chatId, vars.DescriptionMessage)
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: helpers.GetInlineButtonsMain()}

	if _, err := Bot.Send(msg); err != nil {
		errorMsg(vars.HandleDefault, chatId, err)
		log.Println(err)
	}
}
