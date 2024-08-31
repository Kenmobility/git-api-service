package domains

import (
	"time"

	"github.com/kenmobility/git-api-service/internal/http/dtos"
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

func (c Commit) ToDto() dtos.CommitResponseDto {
	return dtos.CommitResponseDto{
		CommitID:   c.CommitID,
		Message:    c.Message,
		Author:     c.Author,
		Date:       c.Date,
		URL:        c.URL,
		Repository: c.RepositoryName,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}
