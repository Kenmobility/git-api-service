package persistence

import (
	"time"

	"github.com/kenmobility/github-api-service/internal/domains/models"
)

// Commit represents the GORM model for the commits table.
type Commit struct {
	ID        uint   `gorm:"primaryKey"`
	CommitID  string `gorm:"type:varchar(100);uniqueIndex"`
	Message   string `gorm:"type:varchar"`
	Author    string `gorm:"type:varchar"`
	Date      time.Time
	URL       string `gorm:"type:varchar"`
	RepoName  string `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ToDomain converts a PostgresCommit to a generic domain Commit.
func (pc *Commit) ToDomain() *models.Commit {
	return &models.Commit{
		CommitID:       pc.CommitID,
		Message:        pc.Message,
		Author:         pc.Author,
		Date:           pc.Date,
		URL:            pc.URL,
		RepositoryName: pc.RepoName,
	}
}

// FromDomain creates a PostgresCommit from a generic domain Commit.
func FromDomain(c *models.Commit) *Commit {
	return &Commit{
		CommitID: c.CommitID,
		Message:  c.Message,
		Author:   c.Author,
		Date:     c.Date,
		URL:      c.URL,
		RepoName: c.RepositoryName,
	}
}
