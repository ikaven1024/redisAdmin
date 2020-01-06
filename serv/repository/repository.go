package repository

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func Open() (*gorm.DB, error) {
	return gorm.Open("sqlite3", "data/admin.db")
}
