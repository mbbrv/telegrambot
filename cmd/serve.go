// Package cmd /*
package cmd

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
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

			if update.Message == nil {
				continue
			}

			if user, ok, err := mysql.IsAuth(Db, update.Message.Chat.UserName); ok {
				if !update.Message.IsCommand() {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, vars.HandleKeyboard)
					if _, err := bot.Send(msg); err != nil {
						errorMsg(update.Message.Chat.ID, bot, err)
						log.Println(err)

						continue
					}

					continue
				}

				if update.Message.Text == "/care" {
					err := mysql.ChangeCareStatus(Db, &user)
					if err != nil {
						errorMsg(update.Message.Chat.ID, bot, err)
						log.Println(err)

						continue
					}

					textMessage := mysql.GetChangeCareStatus(vars.CareDisabled, vars.CareEnabled, &user)

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, textMessage)

					if _, err := bot.Send(msg); err != nil {
						errorMsg(update.Message.Chat.ID, bot, err)
						log.Println(err)

						continue
					}
				}

				if update.Message.Text == "/appointment" {
					textMessage, err := mysql.GetPreparedAppointment(Db, &user)
					if err != nil {
						errorMsg(update.Message.Chat.ID, bot, err)
						log.Println(err)

						continue
					}

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, textMessage)

					if _, err := bot.Send(msg); err != nil {
						errorMsg(update.Message.Chat.ID, bot, err)
						log.Println(err)

						continue
					}
				}

			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, vars.HandleUnauth+err.Error())

				if _, err := bot.Send(msg); err != nil {
					errorMsg(update.Message.Chat.ID, bot, err)
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
