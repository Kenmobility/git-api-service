package models

import (
	"time"
)

type Commit struct {
	CommitID       string    `json:"commit_id"`
	Message        string    `json:"message"`
	Author         string    `json:"author"`
	Date           time.Time `json:"date"`
	URL            string    `json:"url"`
	RepositoryName string    `json:"repo_name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
