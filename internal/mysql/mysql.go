package mysql

import (
	"database/sql"
)

type Config struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
	Host     string `mapstructure:"host"`
}

type User struct {
	id       int
	username string
	care     bool
}

type Appointments struct {
	id          int
	userId      int
	dateTime    sql.NullString
	place       sql.NullString
	description sql.NullString
	cost        sql.NullInt64
}

func GetPreparedAppointment(DB *sql.DB, user *User) (string, error) {
	appointments, err := GetAppointment(DB, user)
	if err != nil {
		return "", err
	}

	return PrepareAppointment(*appointments), nil
}

func GetAppointment(DB *sql.DB, user *User) (*Appointments, error) {
	var appointments Appointments

	row := DB.QueryRow("SELECT * FROM Appointments WHERE user_id = ? ORDER by datetime DESC LIMIT 1", user.id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	if err := row.Scan(&appointments.id, &appointments.userId, &appointments.dateTime, &appointments.place, &appointments.description, &appointments.cost); err != nil {
		return nil, err
	}

	return &appointments, nil
}

func PrepareAppointment(appointments Appointments) string {
	return ""
}

func ChangeCareStatus(DB *sql.DB, user *User) (string, error) {
	return "", nil
}

func IsAuth(DB *sql.DB, userName string) (*User, bool) {
	var user User
	row := DB.QueryRow("SELECT * FROM Users WHERE username = ?", userName)
	if err := row.Scan(&user.id, &user.username); err != nil {
		return nil, false
	}

	return &user, true
}
