package models

import (
	"database/sql"
	"log"
	"math/rand"
	"telegrambot/internal/helpers"
	"telegrambot/internal/vars"
	"time"
)

type Care struct {
	Id          int
	Description sql.NullString
	Time        sql.NullTime
	TimeAt
}

func (user User) GetPreparedCurrentCare() ([]string, error) {
	var caresResult []string

	cares, err := user.getCurrentCare()
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
	if user.MorningCare.Id == care.Id {
		rand.Seed(time.Now().UnixMicro())
		res = vars.MorningGreetings[rand.Intn(len(vars.MorningGreetings))] + user.TelegramUser.FirstName.String + "!\n\n" +
			"Ваши утренние процедуры\n\n"
	} else {
		log.Println(rand.Intn(len(vars.EveningGreetings)))
		rand.Seed(time.Now().UnixMicro())
		res = vars.EveningGreetings[rand.Intn(len(vars.EveningGreetings))] + user.TelegramUser.FirstName.String + "!\n\n" +
			"Ваши вечерние процедуры\n\n"
	}
	return res + care.Description.String
}

func (user User) getCurrentCare() ([]Care, error) {
	var cares []Care

	now := time.Now()
	seconds := helpers.GetTimeDivisionInSeconds()

	past := now.Add(time.Duration(-seconds) * time.Second)
	future := now.Add(time.Duration(seconds) * time.Second)

	if user.MorningCare.Time.Time.After(past) && user.MorningCare.Time.Time.Before(future) && user.MorningCare.Description.Valid {
		cares = append(cares, *user.MorningCare)
	}
	if user.EveningCare.Time.Time.After(past) && user.EveningCare.Time.Time.Before(future) && user.EveningCare.Description.Valid {
		cares = append(cares, *user.MorningCare)
	}

	return cares, nil
}

func (user *User) SetTimeCare(hours string, minutes string, dayTime string) error {
	timeFormat := hours + ":" + minutes + ":00"

	care := user.GetCareByDayTime(dayTime)

	_, err := Db.Exec(`UPDATE cares SET time = ? WHERE id = ?`, timeFormat, care.Id)
	if err != nil {
		return err
	}

	user.TelegramUser.CareEnabled = !user.TelegramUser.CareEnabled
	return nil
}

//func (care Care) GetPhotoDictionary() string {
//	return path.Join(helpers.GetPhotoDictionary(), care.PhotoDictionary.String)
//}
