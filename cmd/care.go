// Package cmd /*
package cmd

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
	"log"
	"telegrambot/internal/user"
	"telegrambot/internal/vars"
	"time"
)

// careCmd represents the care command
var careCmd = &cobra.Command{
	Use:   "care",
	Short: "Запуск бота для отправки сообщений",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		for {
			processCare()
			time.Sleep(time.Minute * vars.TimeToSleep)
		}
	},
}

func init() {
	rootCmd.AddCommand(careCmd)
}

func processCare() {
	users, err := user.GetAllUsersWithCare()
	if err != nil {
		return
	}

	for _, u := range users {
		cares, err := u.GetPreparedCurrentCare()
		if err != nil {
			log.Println(err)
		}
		for _, care := range cares {
			msg := tgbotapi.NewMessage(u.TelegramUser.ChatId.Int64, care)
			msg.ParseMode = "HTML"
			if _, err := Service.Bot.Send(msg); err != nil {
				log.Println(err)
			}
		}
	}
}
