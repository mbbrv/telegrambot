package mysql

import (
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"telegrambot/internal/care/helpers"
	"time"
)

type User struct {
	Id          int
	Username    string
	Care        bool
	PhoneNumber sql.NullString
	TelegramId  sql.NullInt64
	ChatId      sql.NullInt64
}

type Appointments struct {
	id          int
	userId      int
	dateTime    sql.NullString
	place       sql.NullString
	description sql.NullString
	cost        sql.NullInt64
}

type Care struct {
	id              int
	userId          int
	name            sql.NullString
	description     sql.NullString
	url             sql.NullString
	photoDictionary sql.NullString
	time            sql.NullString
	typeCare        sql.NullString
}

func (user User) GetPreparedAppointment(DB *sql.DB) (string, error) {
	appointments, err := user.getAppointment(DB)
	if err != nil {
		return "", err
	}

	return prepareAppointment(*appointments), nil
}

func (user User) getAppointment(DB *sql.DB) (*Appointments, error) {
	var appointments Appointments

	row := DB.QueryRow(`SELECT * FROM Appointments WHERE user_id = ? ORDER by datetime DESC LIMIT 1`, user.Id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	if err := row.Scan(
		&appointments.id,
		&appointments.userId,
		&appointments.dateTime,
		&appointments.place,
		&appointments.description,
		&appointments.cost,
	); err != nil {
		return nil, err
	}

	return &appointments, nil
}

func prepareAppointment(appointments Appointments) string {
	return strconv.FormatInt(appointments.cost.Int64, 10)
}

func (user *User) ChangeCareStatus(DB *sql.DB) error {
	_, err := DB.Exec(`UPDATE Users SET care = ? WHERE id = ?`, !user.Care, user.Id)
	if err != nil {
		return err
	}

	user.Care = !user.Care
	return nil
}

func (user *User) GetChangeCareStatus(disabled string, enabled string) string {
	if user.Care {
		return enabled
	}

	return disabled
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

	_, err := DB.Exec(`UPDATE Users SET username = ?, telegram_id = ?, chat_id = ?  WHERE id = ?`, message.From.UserName, message.From.ID, message.Chat.ID, id)
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

	_, err := DB.Exec(`UPDATE Users SET telegram_id = ?, chat_id = ?  WHERE id = ?`, message.From.ID, message.Chat.ID, id)
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

func (user User) GetPreparedCurrentCare(DB *sql.DB) ([]string, error) {
	var caresResult []string

	cares, err := user.getCurrentCare(DB)
	if err != nil {
		return nil, err
	}

	for _, care := range cares {
		caresResult = append(caresResult, prepareCurrentCare(care))
	}

	return caresResult, nil
}

func (user User) getCurrentCare(DB *sql.DB) ([]Care, error) {
	var cares []Care

	now := time.Now()
	seconds := helpers.GetTimeDivisionInSeconds()
	layout := "15:04:05"

	past := now.Add(time.Duration(-seconds) * time.Second).Format(layout)
	future := now.Add(time.Duration(seconds) * time.Second).Format(layout)

	rows, err := DB.Query(`SELECT * from Care WHERE user_id = ? AND Care.time BETWEEN ? AND ?`, user.Id, past, future)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		log.Println("ya tut")
		var care Care
		if err = rows.Scan(
			&care.id,
			&care.userId,
			&care.name,
			&care.description,
			&care.url,
			&care.photoDictionary,
			&care.time,
			&care.typeCare,
		); err != nil {
			return nil, err
		}
		cares = append(cares, care)
	}
	err = rows.Close()
	if err != nil {
		return nil, err
	}
	return cares, nil
}

func prepareCurrentCare(care Care) string {
	return care.description.String
}
