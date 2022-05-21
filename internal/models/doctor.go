package models

import (
	"database/sql"
)

type Doctor struct {
	Id               int
	Name             sql.NullString
	PhoneNumber      sql.NullString
	TelegramId       sql.NullString
	TelegramUsername sql.NullString
	WhatsAppUrl      sql.NullString
}

func (appointment Appointment) GetDoctor(DB *sql.DB) (Doctor, error) {
	var doctor Doctor

	row := DB.QueryRow(`SELECT * FROM Doctor WHERE id = ?`, appointment.DoctorId)
	if row.Err() != nil {
		return Doctor{}, row.Err()
	}

	if err := row.Scan(
		&doctor.Id,
		&doctor.Name,
		&doctor.PhoneNumber,
		&doctor.TelegramId,
		&doctor.TelegramUsername,
		&doctor.WhatsAppUrl,
	); err != nil {
		return Doctor{}, err
	}

	return doctor, nil
}
