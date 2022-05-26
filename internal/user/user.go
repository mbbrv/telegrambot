package user

import (
	"github.com/nyaruka/phonenumbers"
	"log"
	"math/rand"
	"strconv"
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

func GetUserByPhoneNumber(phone string) (*User, error) {
	userModel, err := models.GetUserByPhoneNum(phone)
	if err != nil {
		return nil, err
	}

	return &User{
		User:         userModel,
		EveningCare:  models.GetCare(userModel.EveningCareId),
		MorningCare:  models.GetCare(userModel.MorningCareId),
		Appointment:  models.GetAppointmentByPatientId(userModel.Id),
		TelegramUser: models.GetTelegramUser(userModel.TelegramUserId),
	}, nil
}

func GetUserByTgId(id int64) *User {
	userTgModel := models.GetTelegramUserByTgId(id)
	userModel := models.GetUserByTgId(userTgModel.Id)

	return &User{
		User:         userModel,
		EveningCare:  models.GetCare(userModel.EveningCareId),
		MorningCare:  models.GetCare(userModel.MorningCareId),
		Appointment:  models.GetAppointmentByPatientId(userModel.Id),
		TelegramUser: userTgModel,
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

	morningTime, _ := user.MorningCare.Time.Time()
	eveningTime, _ := user.EveningCare.Time.Time()

	if morningTime.After(past) && morningTime.Before(future) && user.MorningCare.Description.Valid {
		cares = append(cares, *user.MorningCare)
	}
	if eveningTime.After(past) && eveningTime.Before(future) && user.EveningCare.Description.Valid {
		cares = append(cares, *user.MorningCare)
	}

	return cares, nil
}

func (user User) GetCareByDayTime(dayTime string) *models.Care {
	if dayTime == "morning" {
		return user.MorningCare
	} else {
		return user.EveningCare
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
	care.Time = []byte(timeFormat)

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

func (user User) PrepareAppointment() string {
	parseTime, _ := time.Parse("2006-01-02T15:04:05Z", user.Appointment.DateTime.String)
	doctor := GetUser(user.Appointment.DoctorId)

	phoneNumber, _ := phonenumbers.Parse(doctor.User.PhoneNumber, "RU")

	var res = "<b>Ваша ближайшая запись:</b>\n\n\n" +
		"🧖‍♀️<b>Процедура:</b> " + user.Appointment.Description.String + "\n\n" +
		"💵<b>Стоимость:</b> " + strconv.FormatInt(user.Appointment.Cost.Int64, 10) + " ₽\n\n" +
		"⏰<b>Дата и время:</b> " + parseTime.Format("15:04 02-01-2006") + "\n\n" +
		"🏥<b>Место:</b> " + user.Appointment.Place.String + "\n\n" +
		"👩🏻‍⚕️<b>Врач:</b> " + doctor.User.Name + "\n\n" +
		"=============================" + "\n\n" +
		"Контакты врача:" + "\n\n" +
		"	📱<b>Telegram:</b> " + doctor.TelegramUser.Username.String + "\n\n" +
		"	📞<b>Экстренная связь:</b> " + phonenumbers.Format(phoneNumber, phonenumbers.INTERNATIONAL)

	return res
}
