package mysql

import (
	"database/sql"
	"strconv"
)

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

	row := DB.QueryRow(`SELECT * FROM Appointments WHERE user_id = ? ORDER by datetime DESC LIMIT 1`, user.id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	if err := row.Scan(&appointments.id, &appointments.userId, &appointments.dateTime, &appointments.place, &appointments.description, &appointments.cost); err != nil {
		return nil, err
	}

	return &appointments, nil
}

func PrepareAppointment(appointments Appointments) string {
	return strconv.FormatInt(appointments.cost.Int64, 10)
}

func ChangeCareStatus(DB *sql.DB, user *User) error {
	_, err := DB.Exec(`UPDATE Users SET care = ? WHERE id = ?`, !user.care, user.id)
	if err != nil {
		return err
	}

	user.care = !user.care
	return nil
}

func GetChangeCareStatus(disabled string, enabled string, user *User) string {
	if user.care {
		return enabled
	}

	return disabled
}

func IsAuth(DB *sql.DB, userName string) (User, bool, error) {
	var user User
	row := DB.QueryRow(`SELECT * FROM Users WHERE username = ?`, userName)
	if err := row.Scan(&user.id, &user.username, &user.care); err != nil {
		return User{}, false, err
	}

	return user, true, nil
}
