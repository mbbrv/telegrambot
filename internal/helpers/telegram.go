package helpers

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegrambot/internal/vars"
)

//GetKeyboardButtonsStart используется
//для получения стартовых кнопок
//аутентификации пользователя.
func GetKeyboardButtonsStart() [][]tgbotapi.KeyboardButton {
	return [][]tgbotapi.KeyboardButton{
		{tgbotapi.NewKeyboardButtonContact(vars.KeyboardButtonMobilePhone)},
		{tgbotapi.NewKeyboardButton(vars.KeyboardButtonUsername)},
	}
}

//GetInlineButtonsMain используется для получения кнопок
// получения записи и отключения/включения отправки ухода.
func GetInlineButtonsMain() [][]tgbotapi.InlineKeyboardButton {
	return [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData(vars.InlineButtonAppointment, "appointment")},
		{tgbotapi.NewInlineKeyboardButtonData(vars.InlineButtonCare, "care")},
	}
}

//GetMessage используется для извлечения сущности Message в апдейте.
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

func GetGreetingsMessage(firstName string) string {
	return firstName + ", \n\n" + vars.GreetingsMessage
}

//GetTimeDivisionInSeconds Половина времени в секундах от времени для сна
func GetTimeDivisionInSeconds() int {
	if vars.TimeToSleep%2 == 1 {
		return vars.TimeToSleep/2*60 + 30
	}
	return vars.TimeToSleep / 2 * 60
}
