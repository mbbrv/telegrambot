package helpers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegrambot/internal/vars"
)

func GetConfigDir() string {
	return "./config"
}

func GetPhotoDictionary() string {
	return "./photos"
}

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
		return int(vars.TimeToSleep/2*60 + 30)
	}
	return int(vars.TimeToSleep / 2 * 60)
}

func GetTimeSetMessage(dayTime string) string {
	var res = ""

	if dayTime == "morning" {
		res = "утренних"
	} else {
		res = "вечерних"
	}
	return vars.DailySetTimeMessage + res + " процедур."
}

func GetTimeChangedSuccessMessage(dayTime string) string {
	var res = ""

	if dayTime == "morning" {
		res = "утренних"
	} else {
		res = "вечерних"
	}

	return vars.TimeChangedSuccessMessage + res + " напоминаний!"
}

func GetPreparedDailyCareMessage(morningTime string, eveningTime string) string {
	return vars.DailyMessage + "\n\n" +
		"<b>Время утренних процедур: " + "</b>" + morningTime + "\n" +
		"<b>Время вечерних процедур: " + "</b>" + eveningTime
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

func GetFrom(update *tgbotapi.Update) *tgbotapi.User {
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From
	} else {
		return update.SentFrom()
	}
}
