package database

import (
	"gorm.io/gorm"
)

type Database interface {
	ConnectDb() (*gorm.DB, error)
	Migrate(db *gorm.DB) error
}
