package persistence

import (
	"context"
	"fmt"

	"github.com/kenmobility/github-api-service/internal/domains/models"
	"github.com/kenmobility/github-api-service/internal/domains/services"
	"github.com/kenmobility/github-api-service/internal/dtos"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type GormCommitRepository struct {
	DB *gorm.DB
}

func NewGormCommitRepository(db *gorm.DB) services.CommitRepository {
	return &GormCommitRepository{
		DB: db,
	}
}

// SaveCommit stores a repository commit into the database
func (gc *GormCommitRepository) SaveCommit(ctx context.Context, commit models.Commit) (*models.Commit, error) {
	err := gc.DB.WithContext(ctx).Create(&commit).Error
	return &commit, err
}

// SaveCommit stores a repository commit into the database
func (gc *GormCommitRepository) SaveCommits(ctx context.Context, commits []models.Commit) error {
	var dbCommits []Commit

	for _, commit := range commits {
		dbCommits = append(dbCommits, *FromDomainCommit(&commit))
	}
	return gc.DB.WithContext(ctx).Create(&dbCommits).Error
}

// GetAllCommitsByRepositoryName fetches all stores commits by repository name
func (gc *GormCommitRepository) AllCommitsByRepository(ctx context.Context, repo models.RepoMetadata, query dtos.APIPagingDto) (*dtos.AllCommitsResponse, error) {
	var dbCommits []Commit

	var count, queryCount int64

	queryInfo, offset := getPaginationInfo(query)

	db := gc.DB.WithContext(ctx).Model(&Commit{}).Where(&Commit{RepositoryName: repo.Name})

	db.Count(&count)

	db = db.Offset(offset).Limit(queryInfo.Limit).
		Order(fmt.Sprintf("commits.%s %s", queryInfo.Sort, queryInfo.Direction)).
		Find(&dbCommits)
	db.Count(&queryCount)

	if db.Error != nil {
		log.Info().Msgf("fetch commits error %v", db.Error.Error())

		return nil, db.Error
	}

	pagingInfo := getPagingInfo(queryInfo, int(count))
	pagingInfo.Count = len(dbCommits)

	return &dtos.AllCommitsResponse{
		Commits:  commitResponse(dbCommits),
		PageInfo: pagingInfo,
	}, nil
}

func (gc *GormCommitRepository) TopCommitAuthorsByRepository(ctx context.Context, repo models.RepoMetadata, limit int) ([]string, error) {
	var authors []string
	err := gc.DB.WithContext(ctx).Model(&models.Commit{}).
		Select("author").
		Where("repository_name = ?", repo.Name).
		Group("author").
		Order("count(author) DESC").
		Limit(limit).
		Find(&authors).Error

	return authors, err
}

func commitResponse(commits []Commit) []models.Commit {
	if len(commits) == 0 {
		return nil
	}

	commitsResponse := make([]models.Commit, 0, len(commits))

	for _, c := range commits {
		cr := models.Commit{
			CommitID:       c.CommitID,
			Message:        c.Message,
			Author:         c.Author,
			Date:           c.Date,
			URL:            c.URL,
			RepositoryName: c.RepositoryName,
			CreatedAt:      c.CreatedAt,
			UpdatedAt:      c.UpdatedAt,
		}

		commitsResponse = append(commitsResponse, cr)
	}

	return commitsResponse
}
