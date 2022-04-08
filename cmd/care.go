// Package cmd /*
package cmd

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
	"log"
	"telegrambot/internal/mysql"
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
			time.Sleep(vars.TimeToSleep * time.Minute)
		}
	},
}

func init() {
	rootCmd.AddCommand(careCmd)
}

func processCare() {
	users, err := mysql.GetAllUsersWithCare(Db)
	if err != nil {
		return
	}

	for _, user := range users {
		cares, err := user.GetPreparedCurrentCare(Db)
		if err != nil {
			log.Println(err)
		}
		for _, care := range cares {
			msg := tgbotapi.NewMessage(user.ChatId.Int64, care)
			if _, err := Bot.Send(msg); err != nil {
				log.Println(err)
			}
		}
	}
}
