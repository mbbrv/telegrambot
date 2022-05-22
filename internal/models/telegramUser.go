package models

import (
	"database/sql"
)

type TelegramUser struct {
	Id                      int64          `db:"id"`
	TelegramId              sql.NullInt64  `db:"telegram_id"`
	ChatId                  sql.NullInt64  `db:"chat_id"`
	IsBot                   bool           `db:"is_bot"`
	FirstName               sql.NullString `db:"first_name"`
	LastName                sql.NullString `db:"last_name"`
	Username                sql.NullString `db:"username"`
	CanJoinGroups           bool           `db:"can_join_groups"`
	CanReadAllGroupMessages bool           `db:"can_read_all_group_messages"`
	SupportsInlineQueries   bool           `db:"supports_inline_queries"`
	LanguageCode            sql.NullString `db:"language_code"`
	CareEnabled             bool           `db:"care_enabled"`
	TimeAt
}

func GetTelegramUser(id int64) *TelegramUser {
	telegramUser := TelegramUser{}
	err := Db.Get(&telegramUser, `SELECT * from telegram_users where id = :id`, id)
	if err != nil {
		return nil
	}

	return &telegramUser
}

func UpdateTelegramUser(careEnabled bool, Id int64) error {
	_, err := Db.Exec(`UPDATE telegram_users SET care_enabled = ? WHERE id = ?`, careEnabled, Id)
	if err != nil {
		return err
	}

	return nil
}
