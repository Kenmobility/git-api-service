package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/kenmobility/git-api-service/common/helpers"
	"github.com/rs/zerolog/log"
	"gopkg.in/go-playground/validator.v9"
)

type Config struct {
	AppEnv                string
	GitHubToken           string `validate:"required"`
	DatabaseHost          string `validate:"required"`
	DatabasePort          string `validate:"required"`
	DatabaseUser          string `validate:"required"`
	DatabasePassword      string `validate:"required"`
	DatabaseName          string `validate:"required"`
	FetchInterval         time.Duration
	GitCommitFetchPerPage int
	GitHubApiBaseURL      string
	DefaultStartDate      time.Time
	DefaultEndDate        time.Time
	DefaultRepository     string `validate:"required"`
	Address               string
	Port                  string
}

func LoadConfig(path string) (*Config, error) {
	var err error

	if path == "" {
		path = ".env"
	}
	if err := godotenv.Load(path); err != nil {
		log.Error().Msgf("env config error: %v", err)
		return nil, err
	}

	interval := os.Getenv("FETCH_INTERVAL")
	if interval == "" {
		interval = "1h"
	}

	intervalDuration, err := time.ParseDuration(interval)
	if err != nil {
		log.Error().Msgf("Invalid FETCH_INTERVAL :[%s] env format: %v", interval, err)
		return nil, err
	}

	var sDate time.Time
	var eDate time.Time

	startDate := os.Getenv("DEFAULT_START_DATE")
	if startDate == "" {
		sDate = time.Now().AddDate(0, -10, 0)
	} else {
		sDate, err = time.Parse(time.RFC3339, startDate)
		if err != nil {
			log.Error().Msgf("Invalid DEFAULT_START_DATE [%s] env format: %v", startDate, err)
			return nil, err
		}
	}

	perPage := os.Getenv("GIT_COMMIT_FETCH_PER_PAGE")
	commitPerPage, err := strconv.Atoi(perPage)
	if err != nil {
		commitPerPage = 50
		log.Error().Msgf("Invalid GIT_COMMIT_FETCH_PER_PAGE [%s] env format passed, setting to 50: %v", perPage, err)
	}

	endDate := os.Getenv("DEFAULT_END_DATE")
	if endDate == "" {
		eDate = time.Now()
	} else {
		eDate, err = time.Parse(time.RFC3339, endDate)
		if err != nil {
			log.Info().Msgf("Invalid DEFAULT_END_DATE [%s] env format: %v", endDate, err)
			return nil, err
		}
	}

	configVar := Config{
		AppEnv:                helpers.Getenv("APP_ENV", "local"),
		GitHubToken:           os.Getenv("GIT_HUB_TOKEN"),
		DatabaseHost:          os.Getenv("DATABASE_HOST"),
		DatabasePort:          os.Getenv("DATABASE_PORT"),
		DatabaseUser:          os.Getenv("DATABASE_USER"),
		DatabaseName:          os.Getenv("DATABASE_NAME"),
		DatabasePassword:      os.Getenv("DATABASE_PASSWORD"),
		FetchInterval:         intervalDuration,
		DefaultStartDate:      sDate,
		DefaultEndDate:        eDate,
		GitCommitFetchPerPage: commitPerPage,
		GitHubApiBaseURL:      os.Getenv("GITHUB_API_BASE_URL"),
		Address:               helpers.Getenv("ADDRESS", "127.0.0.1"),
		Port:                  helpers.Getenv("PORT", "5000"),
		DefaultRepository:     helpers.Getenv("DEFAULT_REPOSITORY", "chromium/chromium"),
	}

	validate := validator.New()
	err = validate.Struct(configVar)
	if err != nil {
		log.Info().Msgf("env validation error: %s", err.Error())
		return nil, err
	}

	return &configVar, nil
}
