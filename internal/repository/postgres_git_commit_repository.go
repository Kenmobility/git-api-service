package repository

import (
	"context"
	"fmt"

	"github.com/kenmobility/git-api-service/internal/domains"
	"github.com/kenmobility/git-api-service/internal/http/dtos"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type PostgresGitCommitRepository struct {
	DB *gorm.DB
}

func NewPostgresGitCommitRepository(db *gorm.DB) CommitRepository {
	return &PostgresGitCommitRepository{
		DB: db,
	}
}

// SaveCommit stores a repository commit into the database
func (gc *PostgresGitCommitRepository) SaveCommit(ctx context.Context, commit domains.Commit) (*domains.Commit, error) {
	err := gc.DB.WithContext(ctx).Create(&commit).Error
	return &commit, err
}

// GetAllCommitsByRepositoryName fetches all stores commits by repository name
func (gc *PostgresGitCommitRepository) AllCommitsByRepository(ctx context.Context, repo domains.RepoMetadata, query dtos.APIPagingDto) (*dtos.AllCommitsResponse, error) {
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

func (gc *PostgresGitCommitRepository) TopCommitAuthorsByRepository(ctx context.Context, repo domains.RepoMetadata, limit int) ([]string, error) {
	var authors []string
	err := gc.DB.WithContext(ctx).Model(&domains.Commit{}).
		Select("author").
		Where("repository_name = ?", repo.Name).
		Group("author").
		Order("count(author) DESC").
		Limit(limit).
		Find(&authors).Error

	return authors, err
}

func commitResponse(commits []Commit) []dtos.CommitResponseDto {
	if len(commits) == 0 {
		return nil
	}

	commitsResponse := make([]dtos.CommitResponseDto, 0, len(commits))

	for _, c := range commits {
		cr := dtos.CommitResponseDto{
			CommitID:   c.CommitID,
			Message:    c.Message,
			Author:     c.Author,
			Date:       c.Date,
			URL:        c.URL,
			Repository: c.RepositoryName,
			CreatedAt:  c.CreatedAt,
			UpdatedAt:  c.UpdatedAt,
		}

		commitsResponse = append(commitsResponse, cr)
	}

	return commitsResponse
}
