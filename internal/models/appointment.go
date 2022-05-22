package models

import (
	"database/sql"
	"log"
)

type Appointment struct {
	Id          int64          `db:"id"`
	PatientId   int64          `db:"patient_id"`
	DoctorId    int64          `db:"doctor_id"`
	DateTime    sql.NullString `db:"date_time"`
	Place       sql.NullString `db:"place"`
	Type        sql.NullString `db:"type"`
	Description sql.NullString `db:"description"`
	Cost        sql.NullInt64  `db:"cost"`
	Active      int            `db:"active"`
	TimeAt
}

func GetAppointment(id int64) *Appointment {
	appointment := Appointment{}

	err := Db.Get(&appointment, `SELECT * FROM appointments WHERE id=?`, id)
	if err != nil {
		return nil
	}

	return &appointment
}

func GetAppointmentByPatientId(userId int64) *Appointment {
	appointment := Appointment{}

	err := Db.Get(&appointment, `SELECT * FROM appointments WHERE patient_id=? AND active = 1`, userId)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &appointment
}
