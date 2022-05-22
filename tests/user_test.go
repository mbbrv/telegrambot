package tests

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"log"
	"telegrambot/internal/helpers"
	"telegrambot/internal/models"
	mysql "telegrambot/internal/mysql/config"
	telegram "telegrambot/internal/telegram/config"
	"telegrambot/internal/user"
	"testing"
)

func TestUserModel(t *testing.T) {
	db := GetDbConnection()

	user := models.User{}

	err := db.Get(&user, "SELECT * from users")
	if err != nil {
		log.Println(err)
	}

	b, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(b))
}

func TestUserPhoneNumber(t *testing.T) {
	db := GetDbConnection()
	models.Db = db
	var phone = "79169641992"
	u := models.User{}

	err := db.Get(&u, "SELECT * FROM users where phone_number=?", phone)
	if err != nil {
		log.Println(err)
	}

	b, err := json.Marshal(&u)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(b))

	b, err = json.Marshal(&user.User{
		User:         &u,
		EveningCare:  models.GetCare(u.EveningCareId),
		MorningCare:  models.GetCare(u.MorningCareId),
		Appointment:  models.GetAppointmentByPatientId(u.Id),
		TelegramUser: models.GetTelegramUser(u.TelegramUserId),
	})
	if err != nil {
		log.Println(err)
	}
	log.Println(string(b))
}

func GetDbConnection() *sqlx.DB {
	viper.AddConfigPath(helpers.GetConfigDir())
	viper.SetConfigName("telegramBot")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("start")
	configTg := telegram.Config{}
	err = viper.UnmarshalKey("telegram", &configTg)
	if err != nil {
		log.Fatalln(err)
	}

	var configDB *mysql.Config
	err = viper.UnmarshalKey("mysql", &configDB)
	if err != nil {
		log.Fatalln(err)
	}
	var dataSourceName = configDB.User + ":" + configDB.Password + "@/" + configDB.DB + "?parseTime=true"
	db, err := sqlx.Open("mysql", dataSourceName)
	db.Driver()
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
