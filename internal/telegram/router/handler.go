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

func (r Router) handleChangeTime(data keyboards.Data) (string, error) {
	oldInlineKeyboard := r.update.Message.ReplyMarkup.InlineKeyboard
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
	msg := tgbotapi.NewEditMessageTextAndMarkup(r.update.Message.Chat.ID, r.update.Message.MessageID, helpers.GetTimeSetMessage(data.DayTime), replyMarkup)

	if _, err := r.bot.Send(msg); err != nil {
		log.Println(err)
		return vars.ErrorDefault, err
	}

	return "", nil
}

func (r Router) handleSetTime(dayTime string) (string, error) {
	oldInlineKeyboard := r.update.Message.ReplyMarkup.InlineKeyboard
	hours := oldInlineKeyboard[1][0]
	minutes := oldInlineKeyboard[1][1]
	care := r.user.GetCareByDayTime(dayTime)
	timeCare, _ := time.Parse("15:04:05", string(care.Time))
	timeChanged, _ := time.Parse("15:04:05", hours.Text+":"+minutes.Text+":00")
	log.Println(timeChanged, timeCare)
	if timeChanged != timeCare {
		if err := r.user.SetTimeCare(hours.Text, minutes.Text, dayTime); err != nil {
			return vars.ErrorDefault, err
		}

		msgNew := tgbotapi.NewMessage(r.update.Message.Chat.ID, helpers.GetTimeChangedSuccessMessage(dayTime))
		message, err := r.bot.Send(msgNew)
		if err != nil {
			log.Println(err)
			return vars.ErrorDefault, err
		}
		go func() {
			time.Sleep(3 * time.Second)
			msgDel := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
			if _, err := r.bot.Send(msgDel); err != nil {
				log.Println(err)
				//return vars.ErrorDefault, err
			}
		}()
	}

	preparedDailyMsg := helpers.GetPreparedDailyCareMessage(string(r.user.MorningCare.Time), string(r.user.EveningCare.Time))
	replyMarkup := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboards.GetInlineButtonsDaily("description")}
	msg := tgbotapi.NewEditMessageTextAndMarkup(r.update.Message.Chat.ID, r.update.Message.MessageID, preparedDailyMsg, replyMarkup)
	msg.ParseMode = "HTML"
	if _, err := r.bot.Send(msg); err != nil {
		log.Println(err)
		return vars.ErrorDefault, err
	}

	return "", nil
}

func (r Router) handleTime(dayTime string) (string, error) {
	care := r.user.GetCareByDayTime(dayTime)

	t, err := care.Time.Time()
	if err != nil {
		return vars.ErrorDefault, err
	}

	replyMarkup := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboards.GetInlineButtonsDailyTime(strconv.Itoa(t.Hour()), strconv.Itoa(t.Minute()), dayTime)}
	msg := tgbotapi.NewEditMessageTextAndMarkup(r.update.Message.Chat.ID, r.update.Message.MessageID, helpers.GetTimeSetMessage(dayTime), replyMarkup)

	if _, err := r.bot.Send(msg); err != nil {
		log.Println(err)
		return vars.ErrorDefault, err
	}

	return "", nil
}

func (r Router) handleDaily() (string, error) {
	preparedCareMsg := helpers.GetPreparedDailyCareMessage(string(r.user.MorningCare.Time), string(r.user.EveningCare.Time))
	replyMarkup := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboards.GetInlineButtonsDaily("description")}
	msg := tgbotapi.NewEditMessageTextAndMarkup(r.update.Message.Chat.ID, r.update.Message.MessageID, preparedCareMsg, replyMarkup)
	msg.ParseMode = "HTML"

	if _, err := r.bot.Send(msg); err != nil {
		log.Println(err)
		return vars.ErrorDefault, err
	}

	return "", nil
}

func (r Router) handleCare() (string, error) {
	err := r.user.ChangeCareStatus()
	if err != nil {
		log.Println(err)
		return vars.ErrorDefault, err
	}

	replyMarkup := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboards.GetInlineButtonsMain(r.user.TelegramUser.CareEnabled)}
	msg := tgbotapi.NewEditMessageReplyMarkup(r.update.Message.Chat.ID, r.update.Message.MessageID, replyMarkup)
	//msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboards.GetInlineButtonsMain()}
	if _, err := r.bot.Send(msg); err != nil {
		log.Println(err)
		return vars.ErrorDefault, err
	}

	textMessage := r.user.GetChangeCareStatus(vars.CareDisabled, vars.CareEnabled)
	msgNew := tgbotapi.NewMessage(r.update.Message.Chat.ID, textMessage)
	message, err := r.bot.Send(msgNew)
	if err != nil {
		log.Println(err)
		return vars.ErrorDefault, err
	}
	go func() {
		time.Sleep(3 * time.Second)
		msgDel := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
		if _, err := r.bot.Send(msgDel); err != nil {
			log.Println(err)
		}
	}()

	return "", nil
}

func (r Router) handleAppointment() (string, error) {
	if r.user.Appointment == nil {
		return vars.NoAppointmentMessage, nil
	}
	textMessage := r.user.PrepareAppointment()

	msg := tgbotapi.NewMessage(r.update.Message.Chat.ID, textMessage)
	//msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboards.GetInlineButtonsMain()}
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true
	if _, err := r.bot.Send(msg); err != nil {
		log.Println(err)
		return vars.ErrorDefault, err
	}
	return "", nil
}

func (r Router) HandleGreetings() (string, error) {
	msg := tgbotapi.NewMessage(r.update.Message.Chat.ID, helpers.GetGreetingsMessage(r.user.TelegramUser.FirstName.String))
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

func (r Router) handleDescription() (string, error) {
	msg := tgbotapi.NewMessage(r.update.Message.Chat.ID, vars.DescriptionMessage)
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboards.GetInlineButtonsMain(r.user.TelegramUser.CareEnabled)}

	if _, err := r.bot.Send(msg); err != nil {
		return vars.ErrorDefault, err
	}

	return "", nil
}

func (r Router) handleDescriptionEdit() (string, error) {
	replyMarkup := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboards.GetInlineButtonsMain(r.user.TelegramUser.CareEnabled)}
	msg := tgbotapi.NewEditMessageTextAndMarkup(r.update.Message.Chat.ID, r.update.Message.MessageID, vars.DescriptionMessage, replyMarkup)

	if _, err := r.bot.Send(msg); err != nil {
		return vars.ErrorDefault, err
	}

	return "", nil
}
