package domain

import (
	"time"
)

type Commit struct {
	CommitID       string
	Message        string
	Author         string
	Date           time.Time
	URL            string
	RepositoryName string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type AuthorCommitCount struct {
	Author      string
	CommitCount int
}
