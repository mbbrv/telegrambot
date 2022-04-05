// Package cmd /*
package cmd

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		bot, err := tgbotapi.NewBotAPI(configTelegram.Token)
		if err != nil {
			panic(err)
		}

		bot.Debug = configTelegram.Dev

		updateConfig := tgbotapi.NewUpdate(0)

		updateConfig.Timeout = 30

		updates := bot.GetUpdatesChan(updateConfig)

		for update := range updates {

			message, err := helpers.GetMessage(update)
			if err != nil {
				errorMsg(message.Chat.ID, bot, err)
				log.Println(err)

				continue
			}

			if user, ok, err := mysql.IsAuth(Db, message.Chat.UserName); ok {
				if !message.IsCommand() && update.CallbackQuery == nil {

					//Возможно, ничего не надо отправлять при удалении сообщения
					//msg := tgbotapi.NewMessage(message.Chat.ID, vars.HandleKeyboard)
					//if _, err := bot.Send(msg); err != nil {
					//	errorMsg(message.Chat.ID, bot, err)
					//	log.Println(err)
					//
					//	continue
					//}

					del := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
					if _, err := bot.Send(del); err != nil {
						log.Println(err)

						continue
					}

					continue
				}

				if message.Text == "/start" {
					msg := tgbotapi.NewMessage(message.Chat.ID, vars.WelcomeMessage)
					msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: helpers.GetInlineButtons()}

					if _, err := bot.Send(msg); err != nil {
						errorMsg(update.Message.Chat.ID, bot, err)
						log.Println(err)

						continue
					}
				}

				if update.CallbackData() == "care" {
					err := user.ChangeCareStatus(Db)

					if err != nil {
						errorMsg(message.Chat.ID, bot, err)
						log.Println(err)

						continue
					}

					textMessage := user.GetChangeCareStatus(vars.CareDisabled, vars.CareEnabled)
					msg := tgbotapi.NewMessage(message.Chat.ID, textMessage)

					if _, err := bot.Send(msg); err != nil {
						errorMsg(message.Chat.ID, bot, err)
						log.Println(err)

						continue
					}
				}

				if update.CallbackData() == "appointment" {
					textMessage, err := user.GetPreparedAppointment(Db)
					if err != nil {
						errorMsg(message.Chat.ID, bot, err)
						log.Println(err)

						continue
					}

					msg := tgbotapi.NewMessage(message.Chat.ID, textMessage)
					msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: helpers.GetInlineButtons()}

					if _, err := bot.Send(msg); err != nil {
						errorMsg(message.Chat.ID, bot, err)
						log.Println(err)

						continue
					}
				}

			} else {
				msg := tgbotapi.NewMessage(message.Chat.ID, vars.HandleUnauth+" | "+err.Error())

				if _, err := bot.Send(msg); err != nil {
					errorMsg(message.Chat.ID, bot, err)
					log.Println(err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	err := viper.UnmarshalKey("telegram", &configTelegram)
	if err != nil {
		return
	}
}

func errorMsg(chatId int64, bot *tgbotapi.BotAPI, err error) {
	msg := tgbotapi.NewMessage(chatId, "Ошибка при отправке сообщения(((")
	if configTelegram.Dev {
		msg = tgbotapi.NewMessage(chatId, err.Error())
	}

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}
