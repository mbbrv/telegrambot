package router

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"telegrambot/internal/helpers"
	"telegrambot/internal/telegram/keyboards"
	"telegrambot/internal/vars"
	"time"
)

func (r router) handleChangeTime(data keyboards.Data) (string, error) {
	oldInlineKeyboard := r.message.ReplyMarkup.InlineKeyboard
	hours, err := strconv.Atoi(oldInlineKeyboard[1][0].Text)
	if err != nil {
		log.Println(err)
		return vars.ErrorDefault, err
	}
	minutes, err := strconv.Atoi(oldInlineKeyboard[1][1].Text)
	if err != nil {
		log.Println(err)
		return vars.ErrorDefault, err
	}

	if data.Action == "inc" {
		if data.Data == "hours" {
			hours = helpers.IncreaseHours(hours)
		}
		if data.Data == "minutes" {
			minutes = helpers.IncreaseMinutes(minutes)
		}
	}
	if data.Action == "dec" {
		if data.Data == "hours" {
			hours = helpers.DecreaseHours(hours)
		}
		if data.Data == "minutes" {
			minutes = helpers.DecreaseMinutes(minutes)
		}
	}

	replyMarkup := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboards.GetInlineButtonsDailyTime(strconv.Itoa(hours), strconv.Itoa(minutes), data.DayTime)}
	msg := tgbotapi.NewEditMessageTextAndMarkup(r.message.Chat.ID, r.message.MessageID, helpers.GetTimeSetMessage(data.DayTime), replyMarkup)

	if _, err := r.bot.Send(msg); err != nil {
		log.Println(err)
		return vars.ErrorDefault, err
	}

	return "", nil
}

func (r router) handleSetTime(dayTime string) (string, error) {
	oldInlineKeyboard := r.message.ReplyMarkup.InlineKeyboard
	hours := oldInlineKeyboard[1][0]
	minutes := oldInlineKeyboard[1][1]

	if err := r.user.SetTimeCare(hours.Text, minutes.Text, dayTime, r.db); err != nil {
		return vars.ErrorDefault, err
	}

	replyMarkup := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboards.GetInlineButtonsDaily("description")}
	msg := tgbotapi.NewEditMessageTextAndMarkup(r.message.Chat.ID, r.message.MessageID, vars.DailyMessage, replyMarkup)

	if _, err := r.bot.Send(msg); err != nil {
		log.Println(err)
		return vars.ErrorDefault, err
	}

	return "", nil
}

func (r router) handleTime(dayTime string) (string, error) {
	care, err := r.user.GetCareByDayTime(dayTime, r.db)
	if err != nil {
		return vars.ErrorDefault, err
	}

	t, err := time.Parse("15:04:05", care.Time.String)
	if err != nil {
		return vars.ErrorDefault, err
	}

	replyMarkup := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboards.GetInlineButtonsDailyTime(strconv.Itoa(t.Hour()), strconv.Itoa(t.Minute()), dayTime)}
	msg := tgbotapi.NewEditMessageTextAndMarkup(r.message.Chat.ID, r.message.MessageID, helpers.GetTimeSetMessage(dayTime), replyMarkup)

	if _, err := r.bot.Send(msg); err != nil {
		log.Println(err)
		return vars.ErrorDefault, err
	}

	return "", nil
}

func (r router) handleDaily() (string, error) {
	replyMarkup := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboards.GetInlineButtonsDaily("description")}
	msg := tgbotapi.NewEditMessageTextAndMarkup(r.message.Chat.ID, r.message.MessageID, vars.DailyMessage, replyMarkup)

	if _, err := r.bot.Send(msg); err != nil {
		log.Println(err)
		return vars.ErrorDefault, err
	}

	return "", nil
}

func (r router) handleCare() (string, error) {
	err := r.user.ChangeCareStatus(r.db)
	if err != nil {
		log.Println(err)
		return vars.ErrorDefault, err
	}

	textMessage := r.user.GetChangeCareStatus(vars.CareDisabled, vars.CareEnabled)
	msg := tgbotapi.NewMessage(r.message.Chat.ID, textMessage)
	if _, err := r.bot.Send(msg); err != nil {
		log.Println(err)
		return vars.ErrorDefault, err
	}

	return "", nil
}

func (r router) handleAppointment() (string, error) {
	textMessage, err := r.user.GetPreparedAppointment(r.db)
	if err != nil {
		log.Println(err)
		return vars.ErrorDefault, err
	}

	msg := tgbotapi.NewMessage(r.message.Chat.ID, textMessage)
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboards.GetInlineButtonsMain()}
	if _, err := r.bot.Send(msg); err != nil {
		log.Println(err)
		return vars.ErrorDefault, err
	}

	return "", nil
}

func (r router) handleGreetings() (string, error) {
	msg := tgbotapi.NewMessage(r.message.Chat.ID, helpers.GetGreetingsMessage(r.user.FirstName.String))
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	if _, err := r.bot.Send(msg); err != nil {
		return vars.ErrorDefault, err
	}

	errMsg, err := r.handleDescription()
	if err != nil {
		return errMsg, err
	}
	return "", nil
}

func (r router) handleDescription() (string, error) {
	msg := tgbotapi.NewMessage(r.message.Chat.ID, vars.DescriptionMessage)
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboards.GetInlineButtonsMain()}

	if _, err := r.bot.Send(msg); err != nil {
		return vars.ErrorDefault, err
	}

	return "", nil
}

func (r router) handleDescriptionEdit() (string, error) {
	replyMarkup := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboards.GetInlineButtonsMain()}
	msg := tgbotapi.NewEditMessageTextAndMarkup(r.message.Chat.ID, r.message.MessageID, vars.DescriptionMessage, replyMarkup)

	if _, err := r.bot.Send(msg); err != nil {
		return vars.ErrorDefault, err
	}

	return "", nil
}
