package domain

import (
	"time"
)

type RepoMetadata struct {
	PublicID          string
	Name              string
	Description       string
	URL               string
	Language          string
	ForksCount        int
	StarsCount        int
	OpenIssuesCount   int
	WatchersCount     int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	LastFetchedCommit string
	LastFetchedPage   int32
	IsFetching        bool
}
