package repository

import (
	"time"

	"github.com/kenmobility/git-api-service/internal/domain"
)

// Commit represents the GORM model for the commits table.
type Commit struct {
	ID             uint   `gorm:"primaryKey"`
	CommitID       string `gorm:"type:varchar(100);uniqueIndex"`
	Message        string `gorm:"type:varchar"`
	Author         string `gorm:"type:varchar"`
	Date           time.Time
	URL            string `gorm:"type:varchar"`
	RepositoryName string `gorm:"type:varchar(100);index"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ToDomain converts a PostgresCommit to a generic domain entity Commit.
func (pc *Commit) ToDomain() *domain.Commit {
	return &domain.Commit{
		CommitID:       pc.CommitID,
		Message:        pc.Message,
		Author:         pc.Author,
		Date:           pc.Date,
		URL:            pc.URL,
		RepositoryName: pc.RepositoryName,
	}
}

// FromDomain creates a PostgresCommit from a generic domain entity Commit.
func FromDomainCommit(c *domain.Commit) *Commit {
	return &Commit{
		CommitID:       c.CommitID,
		Message:        c.Message,
		Author:         c.Author,
		Date:           c.Date,
		URL:            c.URL,
		RepositoryName: c.RepositoryName,
	}
}
