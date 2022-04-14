package mysql

import (
	"database/sql"
	"github.com/nyaruka/phonenumbers"
	"strconv"
	"time"
)

type Appointment struct {
	Id          int
	UserId      int
	DoctorId    int
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

	doctor, err := appointment.GetDoctor(DB)
	if err != nil {
		return "", err
	}
	//timeAppointment, _ := time.Parse("2006-01-02T15:04:05Z", appointment.DateTime.String)
	//timeNow.Before()
	return prepareAppointment(&appointment, &doctor), nil
}

func prepareAppointment(appointment *Appointment, doctor *Doctor) string {
	parseTime, _ := time.Parse("2006-01-02T15:04:05Z", appointment.DateTime.String)
	phoneNumber, _ := phonenumbers.Parse(doctor.PhoneNumber.String, "RU")

	var res = "<b>Ваша ближайшая запись:</b>\n\n\n" +
		"🧖‍♀️<b>Процедура:</b> " + appointment.Description.String + "\n\n" +
		"💵<b>Стоимость:</b> " + strconv.FormatInt(appointment.Cost.Int64, 10) + " ₽\n\n" +
		"⏰<b>Дата и время:</b> " + parseTime.Format("15:04 02-01-2006") + "\n\n" +
		"🏥<b>Место:</b> " + appointment.Place.String + "\n\n" +
		"👩🏻‍⚕️<b>Врач:</b> " + doctor.Name.String + "\n\n" +
		"=============================" + "\n\n" +
		"Контакты врача:" + "\n\n" +
		"	📱<b>Telegram:</b> " + doctor.TelegramUsername.String + "\n\n" +
		"	📱<b>WhatsApp:</b> " + doctor.WhatsAppUrl.String + "\n\n" +
		"	📞<b>Экстренная связь:</b> " + phonenumbers.Format(phoneNumber, phonenumbers.INTERNATIONAL)

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
		&appointment.DoctorId,
		&appointment.DateTime,
		&appointment.Place,
		&appointment.Description,
		&appointment.Cost,
	); err != nil {
		return Appointment{}, err
	}

	return appointment, nil
}
