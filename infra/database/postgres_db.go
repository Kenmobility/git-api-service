package database

import (
	"fmt"

	"github.com/kenmobility/git-api-service/infra/config"
	postgreSQL "github.com/kenmobility/git-api-service/internal/repository/postgres"
	"github.com/kenmobility/git-api-service/pkg/helpers"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDatabase struct {
	DSN string
	db  *gorm.DB
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

// ConnectDb establishes a postgreSQL database connection or error if not successful
func (p *PostgresDatabase) ConnectDb() (*gorm.DB, error) {
	var err error
	p.db, err = gorm.Open(postgres.Open(p.DSN), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Info().Msgf("failed to connect to postgres database: %v", err)

		return nil, err
	}
	return p.db, nil
}

// Migrate does db schema migration for PostgreSQL
func (p *PostgresDatabase) Migrate() error {
	// Migrate the schema for PostgreSQL
	return p.db.AutoMigrate(&postgreSQL.Repository{}, &postgreSQL.Commit{})
}
