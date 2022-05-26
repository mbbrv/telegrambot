package models

import (
	"database/sql"
)

type User struct {
	Id                    int64          `db:"id"`
	Name                  string         `db:"name"`
	Surname               sql.NullString `db:"surname"`
	Email                 string         `db:"email"`
	EmailVerifiedAt       sql.NullTime   `db:"email_verified_at"`
	Password              string         `db:"password"`
	RememberToken         sql.NullString `db:"remember_token"`
	TelegramUserId        int64          `db:"telegram_user_id"`
	PhoneNumber           string         `db:"phone_number"`
	PhoneNumberVerifiedAt sql.NullString `db:"phone_number_verified_at"`
	MorningCareId         int64          `db:"morning_care_id"`
	EveningCareId         int64          `db:"evening_care_id"`
	TimeAt
}

func GetUserByPhoneNum(phoneNumber string) (*User, error) {
	user := User{}

	err := Db.Get(&user, "SELECT * FROM users where users.phone_number = ?", phoneNumber)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByTgId(id int64) *User {
	user := User{}

	err := Db.Get(&user, "SELECT * FROM users where users.telegram_user_id = ?", id)
	if err != nil {
		return nil
	}

	return &user
}

func GetUser(id int64) *User {
	user := User{}

	err := Db.Get(&user, "SELECT * from users where users.id = ?", id)
	if err != nil {
		return nil
	}

	return &user
}

func GetAllUsers() ([]User, error) {
	var users []User
	err := Db.Select(&users, `SELECT * from users`)
	if err != nil {
		return nil, err
	}

	return users, nil
}
