package mysql

import (
	"database/sql"
	"log"
	"math/rand"
	"path"
	"strings"
	"telegrambot/internal/helpers"
	"telegrambot/internal/vars"
	"time"
)

type Care struct {
	Id              int
	UserId          int
	Description     sql.NullString
	Url             sql.NullString
	PhotoDictionary sql.NullString
	Time            sql.NullString
	DayTime         sql.NullString
}

type TimeCares struct {
	Morning string
	Evening string
}

func (user User) GetPreparedCurrentCare(DB *sql.DB) ([]string, error) {
	var caresResult []string

	cares, err := user.getCurrentCare(DB)
	if err != nil {
		return nil, err
	}

	for _, care := range cares {
		caresResult = append(caresResult, user.prepareCurrentCare(&care))
	}

	return caresResult, nil
}

func (user User) prepareCurrentCare(care *Care) string {
	var res string
	if care.DayTime.String == "morning" {
		rand.Seed(time.Now().UnixMicro())
		res = vars.MorningGreetings[rand.Intn(len(vars.MorningGreetings))] + user.FirstName.String + "!\n\n" +
			"Ваши утренние процедуры\n\n"
	} else {
		log.Println(rand.Intn(len(vars.EveningGreetings)))
		rand.Seed(time.Now().UnixMicro())
		res = vars.EveningGreetings[rand.Intn(len(vars.EveningGreetings))] + user.FirstName.String + "!\n\n" +
			"Ваши вечерние процедуры\n\n"
	}
	return res + care.Description.String
}

func (user User) getCurrentCare(DB *sql.DB) ([]Care, error) {
	var cares []Care

	now := time.Now()
	seconds := helpers.GetTimeDivisionInSeconds()
	layout := "15:04:05"

	past := now.Add(time.Duration(-seconds) * time.Second).Format(layout)
	future := now.Add(time.Duration(seconds) * time.Second).Format(layout)

	rows, err := DB.Query(`SELECT * from Care WHERE user_id = ? AND Care.time BETWEEN ? AND ?`, user.Id, past, future)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var care Care
		if err = rows.Scan(
			&care.Id,
			&care.UserId,
			&care.Description,
			&care.Url,
			&care.PhotoDictionary,
			&care.Time,
			&care.DayTime,
		); err != nil {
			return nil, err
		}
		cares = append(cares, care)
	}
	err = rows.Close()
	if err != nil {
		return nil, err
	}
	return cares, nil
}

func (user *User) GetCareByDayTime(dayTime string, DB *sql.DB) (*Care, error) {
	var care Care

	row := DB.QueryRow(`SELECT * FROM Care WHERE user_id = ? and day_time = ?`, user.Id, dayTime)
	if err := row.Scan(
		&care.Id,
		&care.UserId,
		&care.Description,
		&care.Url,
		&care.PhotoDictionary,
		&care.Time,
		&care.DayTime,
	); err != nil {
		return nil, err
	}

	return &care, nil
}

func (user *User) GetTimeOfCares(DB *sql.DB) (TimeCares, error) {
	var timeCares TimeCares

	rows, err := DB.Query(`SELECT * from Care WHERE user_id = ?`, user.Id)
	if err != nil {
		return TimeCares{}, err
	}
	for rows.Next() {
		var care Care
		if err = rows.Scan(
			&care.Id,
			&care.UserId,
			&care.Description,
			&care.Url,
			&care.PhotoDictionary,
			&care.Time,
			&care.DayTime,
		); err != nil {
			return TimeCares{}, err
		}
		if care.DayTime.String == "morning" {
			morning := strings.Replace(care.Time.String, ":00", "", 1)
			timeCares.Morning = morning
		} else {
			evening := strings.Replace(care.Time.String, ":00", "", 1)
			timeCares.Evening = evening
		}
	}
	return timeCares, nil
}

func (user *User) SetTimeCare(hours string, minutes string, dayTime string, DB *sql.DB) error {
	timeFormat := hours + ":" + minutes + ":00"
	_, err := DB.Exec(`UPDATE Care SET time = ? WHERE user_id = ? AND day_time = ?`, timeFormat, user.Id, dayTime)
	if err != nil {
		return err
	}

	user.Care = !user.Care
	return nil
}

func (care Care) GetPhotoDictionary() string {
	return path.Join(helpers.GetPhotoDictionary(), care.PhotoDictionary.String)
}
