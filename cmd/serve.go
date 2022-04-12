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
	Short: "Запуск бота для взаимодействия с пользователем",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		updateConfig := tgbotapi.NewUpdate(0)

		updateConfig.Timeout = 30

		updates := Bot.GetUpdatesChan(updateConfig)

		for update := range updates {
			if errMsg, err := processServe(update); err != nil {
				errorMsg(errMsg, update.Message.Chat.ID, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func processServe(update tgbotapi.Update) (string, error) {
	var auth = false
	message, err := helpers.GetMessage(update)
	if err != nil {
		return vars.ErrorDefault, err
	}

	if message.Text == "/start" {

		msg := tgbotapi.NewMessage(message.Chat.ID, vars.WelcomeMessage)
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{Keyboard: helpers.GetKeyboardButtonsStart()}

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
		if err := user.UpdateFirstName(Db, message); err != nil {
			return "", err
		}

		if auth {
			err := greetingsMsg(message.Chat.ID, user.FirstName.String)
			if err != nil {
				log.Println(err)
				return vars.ErrorDefault, err
			}
		}

		if update.CallbackData() == "care" {
			err := user.ChangeCareStatus(Db)
			if err != nil {
				log.Println(err)
				return vars.ErrorDefault, err
			}

			textMessage := user.GetChangeCareStatus(vars.CareDisabled, vars.CareEnabled)
			msg := tgbotapi.NewMessage(message.Chat.ID, textMessage)
			if _, err := Bot.Send(msg); err != nil {
				log.Println(err)
				return vars.ErrorDefault, err
			}
		}

		if update.CallbackData() == "appointment" {
			textMessage, err := user.GetPreparedAppointment(Db)
			if err != nil {
				log.Println(err)
				return vars.ErrorDefault, err
			}

			msg := tgbotapi.NewMessage(message.Chat.ID, textMessage)
			msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: helpers.GetInlineButtonsMain()}
			if _, err := Bot.Send(msg); err != nil {
				log.Println(err)
				return vars.ErrorDefault, err
			}
		}

		if message.Command() == "description" {
			err := descriptionMsg(message.Chat.ID)
			if err != nil {
				log.Println(err)
				return vars.ErrorDefault, err
			}
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

func greetingsMsg(chatId int64, firstName string) error {
	msg := tgbotapi.NewMessage(chatId, helpers.GetGreetingsMessage(firstName))
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	if _, err := Bot.Send(msg); err != nil {
		return err
	}

	err := descriptionMsg(chatId)
	if err != nil {
		return err
	}
	return nil
}

func descriptionMsg(chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, vars.DescriptionMessage)
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: helpers.GetInlineButtonsMain()}

	if _, err := Bot.Send(msg); err != nil {
		return err
	}

	return nil
}
