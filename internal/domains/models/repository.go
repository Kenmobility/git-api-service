package models

import (
	"time"
)

type RepoMetadata struct {
	PublicID        string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	URL             string    `json:"url"`
	Language        string    `json:"language"`
	ForksCount      int       `json:"forks_count"`
	StarsCount      int       `json:"stars_count"`
	OpenIssuesCount int       `json:"open_issues_count"`
	WatchersCount   int       `json:"watchers_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
