package mysql

import (
	"database/sql"
	"time"
)

type Appointment struct {
	Id          int
	UserId      int
	DateTime    sql.NullString
	Place       sql.NullString
	Description sql.NullString
	Cost        sql.NullInt64
}

func (user User) GetPreparedAppointment(DB *sql.DB) (string, error) {
	appointment, err := user.getAppointment(DB)
	if err != nil {
		return "", err
	}
	//timeAppointment, _ := time.Parse("2006-01-02T15:04:05Z", appointment.DateTime.String)
	//timeNow.Before()
	return prepareAppointment(&appointment), nil
}

func prepareAppointment(appointment *Appointment) string {
	parseTime, _ := time.Parse("2006-01-02T15:04:05Z", appointment.DateTime.String)

	var res = "<b>Ближайшая запись:</b>\n\n\n" +
		"💉<b>Процедура:</b> " + appointment.Description.String + "\n\n" +
		"⏰<b>Дата и время:</b> " + parseTime.Format("15:04 02/01/2006") + "\n\n" +
		"🏥<b>Место:</b> " + appointment.Place.String + "\n\n" +
		"👩🏻‍⚕️<b>Врач:</b> test" + "\n\n" +
		"📞<b>Контакты для связи:</b> test"

	return res
}

func (user User) getAppointment(DB *sql.DB) (Appointment, error) {
	var appointment Appointment

	row := DB.QueryRow(`SELECT * FROM Appointment WHERE user_id = ? AND datetime >= ? ORDER by datetime DESC LIMIT 1`, user.Id, time.Now())
	if row.Err() != nil {
		return Appointment{}, row.Err()
	}

	if err := row.Scan(
		&appointment.Id,
		&appointment.UserId,
		&appointment.DateTime,
		&appointment.Place,
		&appointment.Description,
		&appointment.Cost,
	); err != nil {
		return Appointment{}, err
	}

	return appointment, nil
}
