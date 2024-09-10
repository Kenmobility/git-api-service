package domain

import (
	"time"

	"github.com/kenmobility/git-api-service/internal/http/dtos"
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

func (r RepoMetadata) ToDto() dtos.GitRepoMetadataResponseDto {
	return dtos.GitRepoMetadataResponseDto{
		Id:              r.PublicID,
		Name:            r.Name,
		Description:     r.Description,
		URL:             r.URL,
		Language:        r.Language,
		ForksCount:      r.ForksCount,
		StarsCount:      r.StarsCount,
		OpenIssuesCount: r.OpenIssuesCount,
		WatchersCount:   r.WatchersCount,
		CreatedAt:       r.CreatedAt.Format(time.RFC850),
		UpdatedAt:       r.UpdatedAt.Format(time.RFC850),
	}
}
