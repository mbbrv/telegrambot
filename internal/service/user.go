package service

import (
	"log"
	"math/rand"
	"telegrambot/internal/helpers"
	"telegrambot/internal/models"
	"telegrambot/internal/vars"
	"time"
)

type User struct {
	User         *models.User
	EveningCare  *models.Care
	MorningCare  *models.Care
	Appointment  *models.Appointment
	TelegramUser *models.TelegramUser
}

func GetUser(id int64) *User {
	userModel := models.GetUser(id)

	return &User{
		User:         userModel,
		EveningCare:  models.GetCare(userModel.EveningCareId),
		MorningCare:  models.GetCare(userModel.MorningCareId),
		Appointment:  models.GetAppointmentByPatientId(userModel.Id),
		TelegramUser: models.GetTelegramUser(userModel.TelegramUserId),
	}
}

func GetUserByPhoneNumber(phone string) *User {
	userModel := models.GetUserByPhoneNum(phone)

	return &User{
		User:         userModel,
		EveningCare:  models.GetCare(userModel.EveningCareId),
		MorningCare:  models.GetCare(userModel.MorningCareId),
		Appointment:  models.GetAppointmentByPatientId(userModel.Id),
		TelegramUser: models.GetTelegramUser(userModel.TelegramUserId),
	}
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

func (user User) prepareCurrentCare(care *models.Care) string {
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

func (user User) getCurrentCare() ([]models.Care, error) {
	var cares []models.Care

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

func (user User) GetCareByDayTime(dayTime string) models.Care {
	if dayTime == "morning" {
		return *user.MorningCare
	} else {
		return *user.EveningCare
	}
}

func (user *User) GetChangeCareStatus(disabled string, enabled string) string {
	if user.TelegramUser.CareEnabled {
		return enabled
	}

	return disabled
}

func (user *User) SetTimeCare(hours string, minutes string, dayTime string) error {
	timeFormat := hours + ":" + minutes + ":00"
	care := user.GetCareByDayTime(dayTime)
	err := models.UpdateCare(timeFormat, care.Id)
	if err != nil {
		return err
	}

	user.TelegramUser.CareEnabled = !user.TelegramUser.CareEnabled
	return nil
}

func (user *User) ChangeCareStatus() error {
	err := models.UpdateTelegramUser(!user.TelegramUser.CareEnabled, user.TelegramUser.Id)
	if err != nil {
		return err
	}

	user.TelegramUser.CareEnabled = !user.TelegramUser.CareEnabled
	return nil
}

func GetAllUsersWithCare() ([]User, error) {
	var users []User
	usersModel, err := models.GetAllUsers()
	if err != nil {
		return nil, err
	}
	for _, user := range usersModel {
		users = append(users, *GetUser(user.Id))
	}
	return users, nil
}
