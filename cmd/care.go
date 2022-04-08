// Package cmd /*
package cmd

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
	"log"
	"telegrambot/internal/care/vars"
	"telegrambot/internal/mysql"
	"time"
)

// careCmd represents the care command
var careCmd = &cobra.Command{
	Use:   "care",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		for {
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
			time.Sleep(vars.TimeToSleep * time.Minute)
		}
	},
}

func init() {
	rootCmd.AddCommand(careCmd)
}
