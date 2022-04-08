package mysql

import (
	"database/sql"
	"strconv"
)

type Appointments struct {
	Id          int
	UserId      int
	DateTime    sql.NullString
	Place       sql.NullString
	Description sql.NullString
	Cost        sql.NullInt64
}

func (user User) GetPreparedAppointment(DB *sql.DB) (string, error) {
	appointments, err := user.getAppointment(DB)
	if err != nil {
		return "", err
	}

	return prepareAppointment(&appointments), nil
}

func prepareAppointment(appointments *Appointments) string {
	return strconv.FormatInt(appointments.Cost.Int64, 10)
}

func (user User) getAppointment(DB *sql.DB) (Appointments, error) {
	var appointments Appointments

	row := DB.QueryRow(`SELECT * FROM Appointments WHERE user_id = ? ORDER by datetime DESC LIMIT 1`, user.Id)
	if row.Err() != nil {
		return Appointments{}, row.Err()
	}

	if err := row.Scan(
		&appointments.Id,
		&appointments.UserId,
		&appointments.DateTime,
		&appointments.Place,
		&appointments.Description,
		&appointments.Cost,
	); err != nil {
		return Appointments{}, err
	}

	return appointments, nil
}
