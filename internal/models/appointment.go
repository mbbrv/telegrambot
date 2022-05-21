package models

import (
	"database/sql"
	"github.com/nyaruka/phonenumbers"
	"log"
	"strconv"
	"time"
)

type Appointment struct {
	Id          int            `db:"id"`
	PatientId   int            `db:"patient_id"`
	DoctorId    int            `db:"doctor_id"`
	DateTime    sql.NullString `db:"date_time"`
	Place       sql.NullString `db:"place"`
	Type        sql.NullString `db:"type"`
	Description sql.NullString `db:"description"`
	Cost        sql.NullInt64  `db:"cost"`
	Active      int            `db:"active"`
	TimeAt
}

func GetAppointment(id int) (Appointment, error) {
	appointment := Appointment{}

	err := Db.Get(&appointment, `SELECT * FROM appointments WHERE id = :id`, id)
	if err != nil {
		return Appointment{}, err
	}

	return appointment, nil
}

func GetAppointmentByPatientId(userId int) (Appointment, error) {
	appointment := Appointment{}

	err := Db.Get(&appointment, `SELECT * FROM appointments WHERE patient_id = :patient_id AND active = 1`, userId)
	if err != nil {
		return Appointment{}, err
	}

	return appointment, nil
}

func (appointment Appointment) PrepareAppointment() string {
	parseTime, _ := time.Parse("2006-01-02T15:04:05Z", appointment.DateTime.String)
	doctor, err := GetUser(appointment.DoctorId)
	if err != nil {
		log.Println(err)
	}
	phoneNumber, _ := phonenumbers.Parse(doctor.PhoneNumber, "RU")

	var res = "<b>Ваша ближайшая запись:</b>\n\n\n" +
		"🧖‍♀️<b>Процедура:</b> " + appointment.Description.String + "\n\n" +
		"💵<b>Стоимость:</b> " + strconv.FormatInt(appointment.Cost.Int64, 10) + " ₽\n\n" +
		"⏰<b>Дата и время:</b> " + parseTime.Format("15:04 02-01-2006") + "\n\n" +
		"🏥<b>Место:</b> " + appointment.Place.String + "\n\n" +
		"👩🏻‍⚕️<b>Врач:</b> " + doctor.Name + "\n\n" +
		"=============================" + "\n\n" +
		"Контакты врача:" + "\n\n" +
		"	📱<b>Telegram:</b> " + doctor.TelegramUser.Username.String + "\n\n" +
		"	📞<b>Экстренная связь:</b> " + phonenumbers.Format(phoneNumber, phonenumbers.INTERNATIONAL)

	return res
}
