package helpers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegrambot/internal/vars"
)

//GetMessage используется для извлечения сущности Message в апдейте.
func GetMessage(update *tgbotapi.Update) *tgbotapi.Message {
	if update.Message != nil {
		return update.Message
	}
	update.CallbackData()
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message
	}

	return nil
}

// GetGreetingsMessage получить сообщение приветствия.
func GetGreetingsMessage(firstName string) string {
	return firstName + ", \n\n" + vars.GreetingsMessage
}

//GetTimeDivisionInSeconds Половина времени в секундах от времени для сна.
func GetTimeDivisionInSeconds() int {
	if vars.TimeToSleep%2 == 1 {
		return vars.TimeToSleep/2*60 + 30
	}
	return vars.TimeToSleep / 2 * 60
}

func GetTimeSetMessage(dayTime string) string {
	return vars.DailySetTimeMessage + dayTime + " процедур."
}

func IncreaseHours(hours int) int {
	if hours < 23 {
		return hours + 1
	} else {
		return 0
	}
}

func DecreaseHours(hours int) int {
	if hours > 0 {
		return hours - 1
	} else {
		return 23
	}
}

func IncreaseMinutes(minutes int) int {
	if minutes < 45 {
		return minutes + 15
	} else {
		return 0
	}
}

func DecreaseMinutes(minutes int) int {
	if minutes > 0 {
		return minutes - 15
	} else {
		return 45
	}
}
