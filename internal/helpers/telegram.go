package helpers

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegrambot/internal/vars"
)

func GetInlineButtons() [][]tgbotapi.InlineKeyboardButton {

	return [][]tgbotapi.InlineKeyboardButton{
		{getInlineButtonAppointment()},
		{getInlineButtonCare()},
	}
}

func getInlineButtonAppointment() tgbotapi.InlineKeyboardButton {
	var appointment = "appointment"

	return tgbotapi.InlineKeyboardButton{Text: vars.InlineButtonAppointment, CallbackData: &appointment}
}

func getInlineButtonCare() tgbotapi.InlineKeyboardButton {
	var care = "care"

	return tgbotapi.InlineKeyboardButton{Text: vars.InlineButtonCare, CallbackData: &care}
}

func GetMessage(update tgbotapi.Update) (*tgbotapi.Message, error) {
	if update.Message != nil {
		return update.Message, nil
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message, nil
	}

	return nil, errors.New("error while getting Message")
}
