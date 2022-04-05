package helpers

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegrambot/internal/vars"
)

func GetKeyboardButtonsStart() [][]tgbotapi.KeyboardButton {
	return [][]tgbotapi.KeyboardButton{
		{tgbotapi.NewKeyboardButtonContact(vars.KeyboardButtonMobilePhone)},
		{tgbotapi.NewKeyboardButton(vars.KeyboardButtonUsername)},
	}
}

func GetInlineButtonsMain() [][]tgbotapi.InlineKeyboardButton {
	return [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData(vars.InlineButtonAppointment, "appointment")},
		{tgbotapi.NewInlineKeyboardButtonData(vars.InlineButtonCare, "care")},
	}
}

func GetMessage(update tgbotapi.Update) (*tgbotapi.Message, error) {
	if update.Message != nil {
		return update.Message, nil
	}
	update.CallbackData()
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message, nil
	}

	return nil, errors.New("error while getting Message")
}
