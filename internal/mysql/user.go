package mysql

import (
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct {
	Id          int
	Username    sql.NullString
	Care        bool
	PhoneNumber sql.NullString
	TelegramId  sql.NullInt64
	ChatId      sql.NullInt64
	FirstName   sql.NullString
}

func IsAuth(DB *sql.DB, from *tgbotapi.Chat) (User, bool, error) {
	var user User

	row := DB.QueryRow(`SELECT * FROM Users WHERE username = ? OR telegram_id = ?`, from.UserName, from.ID)
	if err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Care,
		&user.PhoneNumber,
		&user.TelegramId,
		&user.ChatId,
		&user.FirstName,
	); err != nil {
		return User{}, false, err
	}

	return user, true, nil
}

func UserEnrichmentByPhoneNumb(DB *sql.DB, message *tgbotapi.Message) error {
	var id int

	row := DB.QueryRow(`SELECT id FROM Users WHERE phone_number = ?`, message.Contact.PhoneNumber)
	if err := row.Scan(&id); err != nil {
		return err
	}

	_, err := DB.Exec(`UPDATE Users SET username = ?, telegram_id = ?, chat_id = ?, first_name = ?  WHERE id = ?`, message.From.UserName, message.From.ID, message.Chat.ID, message.From.FirstName, id)
	if err != nil {
		return err
	}

	return nil
}

func UserEnrichmentByUsername(DB *sql.DB, message *tgbotapi.Message) error {
	var id int

	row := DB.QueryRow(`SELECT id FROM Users WHERE username = ?`, message.Chat.UserName)
	if err := row.Scan(&id); err != nil {
		return err
	}

	_, err := DB.Exec(`UPDATE Users SET telegram_id = ?, chat_id = ?, first_name = ?  WHERE id = ?`, message.From.ID, message.Chat.ID, message.From.FirstName, id)
	if err != nil {
		return err
	}

	return nil
}

func GetAllUsersWithCare(DB *sql.DB) ([]User, error) {
	var users []User
	rows, err := DB.Query(`SELECT * FROM Users WHERE care = ?`, 1)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user User
		if err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.Care,
			&user.PhoneNumber,
			&user.TelegramId,
			&user.ChatId,
			&user.FirstName,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	err = rows.Close()
	if err != nil {
		return nil, err
	}

	return users, err
}
