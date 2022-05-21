// Package cmd /*
package cmd

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
	"log"
	"telegrambot/internal/helpers"
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
			Service.TgUpdate = &update
			Service.TgUpdate.Message = helpers.GetMessage(Service.TgUpdate)
			if errMsg, err := Service.Run(); err != nil {
				errorMsg(errMsg, Service.TgUpdate.Message.Chat.ID, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func errorMsg(message string, chatId int64, err error) {
	msg := tgbotapi.NewMessage(chatId, message)
	if Service.Config.Dev {
		msg = tgbotapi.NewMessage(chatId, err.Error())
	}

	if _, err := Service.Bot.Send(msg); err != nil {
		log.Println(err)
	}
}
