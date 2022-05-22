package models

import (
	"database/sql"
)

type Care struct {
	Id          int64
	Description sql.NullString
	Time        sql.NullTime
	TimeAt
}

func GetCare(id int64) *Care {
	var care = Care{}
	err := Db.Get(&care, `SELECT * FROM cares WHERE id = :id`, id)
	if err != nil {
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
