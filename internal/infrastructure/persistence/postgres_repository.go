package persistence

import (
	"github.com/kenmobility/github-api-service/internal/domains/models"

	"time"
)

// Repository represents the GORM model for the repositories table.
type Repository struct {
	ID              uint   `gorm:"primarykey"`
	PublicID        string `gorm:"type:varchar;uniqueIndex"`
	Name            string `gorm:"type:varchar;unique"`
	Description     string `gorm:"type:text"`
	URL             string `gorm:"type:varchar"`
	Language        string `gorm:"type:varchar"`
	ForksCount      int
	StarsCount      int
	OpenIssuesCount int
	WatchersCount   int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ToDomain converts a PostgresRepository to a generic domain entity RepoMetadata.
func (pr *Repository) ToDomain() *models.RepoMetadata {
	return &models.RepoMetadata{
		Name:            pr.Name,
		Description:     pr.Description,
		URL:             pr.URL,
		Language:        pr.Language,
		ForksCount:      pr.ForksCount,
		StarsCount:      pr.StarsCount,
		OpenIssuesCount: pr.OpenIssuesCount,
		WatchersCount:   pr.WatchersCount,
		CreatedAt:       pr.CreatedAt,
		UpdatedAt:       pr.UpdatedAt,
	}
}

// FromDomain creates a PostgresRepository from a generic domain entity RepoMetadata.
func FromDomainRepo(r *models.RepoMetadata) *Repository {
	return &Repository{
		Name:            r.Name,
		Description:     r.Description,
		URL:             r.URL,
		Language:        r.Language,
		ForksCount:      r.ForksCount,
		StarsCount:      r.StarsCount,
		OpenIssuesCount: r.OpenIssuesCount,
		WatchersCount:   r.WatchersCount,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}
