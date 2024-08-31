package database

import (
	"fmt"

	"github.com/kenmobility/github-api-service/common/helpers"
	"github.com/kenmobility/github-api-service/config"
	"github.com/kenmobility/github-api-service/internal/infrastructure/persistence"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDatabase struct {
	DSN string
}

func NewPostgresDatabase(config config.Config) Database {
	conString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s",
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseUser,
		config.DatabaseName,
		config.DatabasePassword,
	)

	if helpers.IsLocal() {
		conString += " sslmode=disable"
	}

	return &PostgresDatabase{DSN: conString}

	/*
		if helpers.IsLocal() {
			conString += " sslmode=disable"
		}

		db, err := gorm.Open(postgres.Open(conString), &gorm.Config{})
		if err != nil {
			log.Printf("failed to connect to postgres database: %v", err)
			return nil, err
		}
		return db, nil
	*/
}

func (p *PostgresDatabase) ConnectDb() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(p.DSN), &gorm.Config{})
	if err != nil {
		log.Info().Msgf("failed to connect to postgres database: %v", err)

		return nil, err
	}
	return db, nil
}

func (p *PostgresDatabase) Migrate(db *gorm.DB) error {
	// Migrate the schema for PostgreSQL
	return db.AutoMigrate(&persistence.Repository{}, &persistence.Commit{})
}
