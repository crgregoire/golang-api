package db

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // Importing this for gorm to designate the db driver
)

//
// Open database connection specific to the environment
//
func Open() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", os.Getenv("DB_USER")+"@"+os.Getenv("DB_HOST")+"/"+os.Getenv("DB_NAME")+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}
	return db, nil
}
