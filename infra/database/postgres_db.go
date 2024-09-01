package database

import (
	"fmt"

	"github.com/kenmobility/git-api-service/common/helpers"
	"github.com/kenmobility/git-api-service/infra/config"
	"github.com/kenmobility/git-api-service/infra/logger"
	"github.com/kenmobility/git-api-service/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var log = logger.New("postgres")

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
	return db.AutoMigrate(&repository.Repository{}, &repository.Commit{})
}
