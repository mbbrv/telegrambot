package service

import (
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"log"
	"telegrambot/internal/helpers"
	telegram "telegrambot/internal/telegram/config"
	"telegrambot/internal/telegram/keyboards"
	router2 "telegrambot/internal/telegram/router"
	"telegrambot/internal/user"
	"telegrambot/internal/vars"
)

type Service struct {
	User     *user.User
	Db       *sqlx.DB
	Bot      *tgbotapi.BotAPI
	Config   *telegram.Config
	Auth     bool
	TgUpdate *tgbotapi.Update
}

func (service *Service) Run() (string, error) {
	service.Auth = false
	//if service.TgUpdate.Message == nil {
	//	return "", nil
	//}

	if service.TgUpdate.Message.Text == "/start" {

		msg := tgbotapi.NewMessage(service.TgUpdate.Message.Chat.ID, vars.WelcomeMessage)
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{Keyboard: keyboards.GetKeyboardButtonsStart()}
		msg.ParseMode = "HTML"

		if _, err := service.Bot.Send(msg); err != nil {
			log.Println(err)
			return vars.ErrorDefault, err
		}

		return "", nil
	}

	if !service.TgUpdate.Message.IsCommand() && service.TgUpdate.CallbackQuery == nil && service.TgUpdate.Message.Contact == nil && service.TgUpdate.Message.Text != vars.KeyboardButtonUsername {
		del := tgbotapi.NewDeleteMessage(service.TgUpdate.Message.Chat.ID, service.TgUpdate.Message.MessageID)
		if _, err := service.Bot.Send(del); err != nil {
			log.Println(err)

			return "", nil
		}

		return "", nil
	}

	if service.TgUpdate.Message.Contact != nil {
		err := service.Update()
		if err != nil {
			log.Println(err)
			return vars.ErrorNoUser, err
		}
	} else {
		service.User = user.GetUserByTgId(helpers.GetFrom(service.TgUpdate).ID)
	}

	router := router2.NewRouter(service.User, service.TgUpdate, service.Bot, service.Db, &service.Auth)
	if errMsg, err := router.Route(); err != nil {
		log.Println(err)
		return errMsg, err
	}

	return "", nil
}

func (service *Service) Update() error {
	tgUser := service.TgUpdate.SentFrom()
	userU := user.GetUserByPhoneNumber(service.TgUpdate.Message.Contact.PhoneNumber)

	userU.TelegramUser.TelegramId = sql.NullInt64{Int64: tgUser.ID}
	userU.TelegramUser.ChatId = sql.NullInt64{Int64: service.TgUpdate.Message.Chat.ID}
	userU.TelegramUser.IsBot = tgUser.IsBot
	userU.TelegramUser.FirstName = sql.NullString{String: tgUser.FirstName}
	userU.TelegramUser.LastName = sql.NullString{String: tgUser.LastName}
	userU.TelegramUser.Username = sql.NullString{String: tgUser.UserName}
	userU.TelegramUser.CanJoinGroups = tgUser.CanJoinGroups
	userU.TelegramUser.CanReadAllGroupMessages = tgUser.CanReadAllGroupMessages
	userU.TelegramUser.SupportsInlineQueries = tgUser.SupportsInlineQueries
	userU.TelegramUser.LanguageCode = sql.NullString{String: tgUser.LanguageCode}

	_, err := service.Db.NamedExec(`UPDATE telegram_users SET 
                          telegram_id = :telegram_id, 
                          chat_id = :chat_id, 
                          is_bot = :is_bot, 
                          first_name = :first_name, 
                          last_name = :last_name, 
                          username = :username, 
                          can_join_groups = :cjg,
                          can_read_all_group_messages = :cralgm,
                          supports_inline_queries = :siq,
                          language_code = :lc
                          WHERE id = :id`,
		map[string]interface{}{
			"id":          userU.TelegramUser.Id,
			"telegram_id": tgUser.ID,
			"chat_id":     service.TgUpdate.Message.Chat.ID,
			"is_bot":      tgUser.IsBot,
			"first_name":  tgUser.FirstName,
			"last_name":   tgUser.LastName,
			"username":    tgUser.UserName,
			"cjg":         tgUser.CanJoinGroups,
			"cralgm":      tgUser.CanReadAllGroupMessages,
			"siq":         tgUser.SupportsInlineQueries,
			"lc":          tgUser.LanguageCode,
		})
	if err != nil {
		return err
	}

	service.User = userU
	service.Auth = true

	return nil
}
