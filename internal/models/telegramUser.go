package models

import "database/sql"

type TelegramUser struct {
	Id                      int            `db:"id"`
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
