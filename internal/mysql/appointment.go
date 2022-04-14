package mysql

import (
	"database/sql"
	"time"
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
	//timeAppointment, _ := time.Parse("2006-01-02T15:04:05Z", appointments.DateTime.String)
	//timeNow.Before()
	return prepareAppointment(&appointments), nil
}

func prepareAppointment(appointments *Appointments) string {
	parseTime, _ := time.Parse("2006-01-02T15:04:05Z", appointments.DateTime.String)

	var res = "<b>Ближайшая запись:</b>\n\n\n" +
		"💉<b>Процедура:</b> " + appointments.Description.String + "\n\n" +
		"⏰<b>Дата и время:</b> " + parseTime.Format("15:04 02/01/2006") + "\n\n" +
		"🏥<b>Место:</b> " + appointments.Place.String + "\n\n" +
		"👩🏻‍⚕️<b>Врач:</b> test" + "\n\n" +
		"📞<b>Контакты для связи:</b> test"

	return res
}

func (user User) getAppointment(DB *sql.DB) (Appointments, error) {
	var appointments Appointments

	row := DB.QueryRow(`SELECT * FROM Appointments WHERE user_id = ? AND datetime >= ? ORDER by datetime DESC LIMIT 1`, user.Id, time.Now())
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
