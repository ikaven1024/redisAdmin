package repository

import (
	//"github.com/boltdb/bolt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	//"time"
)

func Open() (*gorm.DB, error) {
	return gorm.Open("sqlite3", "admin.db")
}
