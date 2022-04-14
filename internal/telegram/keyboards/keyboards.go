package keyboards

import (
	"encoding/json"
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
		{tgbotapi.NewInlineKeyboardButtonData(vars.InlineButtonCare, "care"), tgbotapi.NewInlineKeyboardButtonData(vars.InlineButtonDailyCare, "daily")},
	}
}

// GetInlineButtonsDaily кнопки выбора настроек для утреннего и вечернего типов ухода
func GetInlineButtonsDaily(parentInline string) [][]tgbotapi.InlineKeyboardButton {
	return [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData(vars.InlineButtonMorning, "morningTime"), tgbotapi.NewInlineKeyboardButtonData(vars.InlineButtonEvening, "eveningTime")},
		{tgbotapi.NewInlineKeyboardButtonData(vars.InlineButtonBack, parentInline)},
	}
}

type Data struct {
	Action  string `json:"action"`
	Data    string `json:"data"`
	DayTime string `json:"dayTime"`
}

// GetInlineButtonsDailyTime кнопки выбора времени для определенного типа ухода
func GetInlineButtonsDailyTime(hours string, minutes string, dayTime string) [][]tgbotapi.InlineKeyboardButton {
	incHours, _ := json.Marshal(Data{Action: "inc", Data: "hours", DayTime: dayTime})
	incMinutes, _ := json.Marshal(Data{Action: "inc", Data: "minutes", DayTime: dayTime})
	decHours, _ := json.Marshal(Data{Action: "dec", Data: "hours", DayTime: dayTime})
	decMinutes, _ := json.Marshal(Data{Action: "dec", Data: "minutes", DayTime: dayTime})

	return [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData(vars.InlineButtonInc, string(incHours)), tgbotapi.NewInlineKeyboardButtonData(vars.InlineButtonInc, string(incMinutes))},
		{tgbotapi.NewInlineKeyboardButtonData(hours, hours), tgbotapi.NewInlineKeyboardButtonData(minutes, minutes)},
		{tgbotapi.NewInlineKeyboardButtonData(vars.InlineButtonDec, string(decHours)), tgbotapi.NewInlineKeyboardButtonData(vars.InlineButtonDec, string(decMinutes))},
		{tgbotapi.NewInlineKeyboardButtonData(vars.InlineButtonOk, dayTime+"SetTime")},
	}
}
