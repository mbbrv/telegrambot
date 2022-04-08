package mysql

import (
	"database/sql"
	"log"
	"telegrambot/internal/helpers"
	"time"
)

type Care struct {
	Id              int
	UserId          int
	Name            sql.NullString
	Description     sql.NullString
	Url             sql.NullString
	PhotoDictionary sql.NullString
	Time            sql.NullString
	TypeCare        sql.NullString
}

func (user User) GetPreparedCurrentCare(DB *sql.DB) ([]string, error) {
	var caresResult []string

	cares, err := user.getCurrentCare(DB)
	if err != nil {
		return nil, err
	}

	for _, care := range cares {
		caresResult = append(caresResult, prepareCurrentCare(&care))
	}

	return caresResult, nil
}

func prepareCurrentCare(care *Care) string {
	return care.Description.String
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
		log.Println("ya tut")
		var care Care
		if err = rows.Scan(
			&care.Id,
			&care.UserId,
			&care.Name,
			&care.Description,
			&care.Url,
			&care.PhotoDictionary,
			&care.Time,
			&care.TypeCare,
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

func (user *User) ChangeCareStatus(DB *sql.DB) error {
	_, err := DB.Exec(`UPDATE Users SET care = ? WHERE id = ?`, !user.Care, user.Id)
	if err != nil {
		return err
	}

	user.Care = !user.Care
	return nil
}

func (user *User) GetChangeCareStatus(disabled string, enabled string) string {
	if user.Care {
		return enabled
	}

	return disabled
}
