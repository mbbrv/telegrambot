package models

import (
	"database/sql"
)

type User struct {
	Id                    int            `db:"id"`
	Name                  string         `db:"name,omitempty"`
	Surname               sql.NullString `db:"surname"`
	Email                 string         `db:"email,omitempty"`
	EmailVerifiedAt       sql.NullTime   `db:"email_verified_at"`
	TelegramUser          *TelegramUser  `db:"telegram_user"`
	PhoneNumber           string         `db:"phone_number"`
	PhoneNumberVerifiedAt sql.NullString `db:"phone_number_verified_at"`
	MorningCare           *Care          `db:"morning_care"`
	EveningCare           *Care          `db:"evening_care"`
	Appointment           *Appointment   `db:"appointment"`
	TimeAt
}

func GetUserByPhoneNum(phoneNumber string) (User, error) {
	user := User{}

	err := Db.Get(&user, "SELECT * from users INNER JOIN telegram_users ON users.id = telegram_users.id INNER JOIN cares c on users.morning_care_id = c.id INNER JOIN cares cc on users.evening_care_id = cc.id where users.phone_number = :phoneNumber", phoneNumber)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func GetUser(id int) (User, error) {
	user := User{}

	err := Db.Get(&user, "SELECT * from users INNER JOIN telegram_users ON users.id = telegram_users.id INNER JOIN cares c on users.morning_care_id = c.id INNER JOIN cares cc on users.evening_care_id = cc.id where users.id = :id", id)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (user User) GetCareByDayTime(dayTime string) Care {
	if dayTime == "morning" {
		return *user.MorningCare
	} else {
		return *user.EveningCare
	}
}

func GetAllUsers() ([]User, error) {
	var users []User
	err := Db.Select(&users, `SELECT * from users INNER JOIN telegram_users ON users.id = telegram_users.id INNER JOIN cares c on users.morning_care_id = c.id INNER JOIN cares cc on users.evening_care_id = cc.id`)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetAllUsersWithCare() ([]User, error) {
	var users []User
	err := Db.Select(&users, `SELECT * from users INNER JOIN telegram_users ON users.id = telegram_users.id INNER JOIN cares c on users.morning_care_id = c.id INNER JOIN cares cc on users.evening_care_id = cc.id where telegram_users.care_enabled = 1`)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (user *User) ChangeCareStatus() error {
	_, err := Db.Exec(`UPDATE telegram_users SET care_enabled = ? WHERE id = ?`, !user.TelegramUser.CareEnabled, user.TelegramUser.Id)
	if err != nil {
		return err
	}

	user.TelegramUser.CareEnabled = !user.TelegramUser.CareEnabled
	return nil
}

func (user *User) GetChangeCareStatus(disabled string, enabled string) string {
	if user.TelegramUser.CareEnabled {
		return enabled
	}

	return disabled
}
