package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/kenmobility/git-api-service/common/helpers"
	"github.com/kenmobility/git-api-service/infra/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func TestMain(m *testing.M) {
	var err error
	config, err := config.LoadConfig("../.env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

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

	db, err = gorm.Open(postgres.Open(conString), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to postgres database: %v", err)
	}

	os.Exit(m.Run())
}
