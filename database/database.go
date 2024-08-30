package database

import (
	"gorm.io/gorm"
)

type Database interface {
	ConnectDb() (*gorm.DB, error)
	Migrate(db *gorm.DB) error
}

/* NewDatabase creates a connection to db and returns the db instance
func NewDatabase(config config.Config) (*Database, error) {
	postgreDb, err := connectPostgresDb(config)
	if err != nil {
		return nil, err
	}
	return &Database{
		Db: postgreDb,
	}, nil
}
*/
