package models

import (
	"database/sql"
	"log"
)

type Care struct {
	Id          int64          `db:"id"`
	Description sql.NullString `db:"description"`
	Time        rawTime        `db:"time"`
	TimeAt
}

func GetCare(id int64) *Care {
	var care = Care{}
	err := Db.Get(&care, `SELECT * FROM cares WHERE id=?`, id)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &care
}

func UpdateCare(timeFormat string, Id int64) error {
	_, err := Db.Exec(`UPDATE cares SET time = ? WHERE id = ?`, timeFormat, Id)
	if err != nil {
		return err
	}
	return nil
}

//func (care Care) GetPhotoDictionary() string {
//	return path.Join(helpers.GetPhotoDictionary(), care.PhotoDictionary.String)
//}
