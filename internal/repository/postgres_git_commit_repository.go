package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/kenmobility/git-api-service/common/message"
	"github.com/kenmobility/git-api-service/internal/domain"
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

// GetByCommitID fetches a commit using commit ID
func (gc *PostgresGitCommitRepository) GetByCommitID(ctx context.Context, commitID string) (*domain.Commit, error) {
	if ctx.Err() == context.Canceled {
		return nil, message.ErrContextCancelled
	}
	var commit Commit
	err := gc.DB.WithContext(ctx).Where("commit_id = ?", commitID).Find(&commit).Error

	if commit.ID == 0 {
		return nil, message.ErrNoRecordFound
	}
	return commit.ToDomain(), err
}

// SaveCommit stores a repository commit into the database
func (gc *PostgresGitCommitRepository) SaveCommit(ctx context.Context, commit domain.Commit) (*domain.Commit, error) {
	if ctx.Err() == context.Canceled {
		return nil, message.ErrContextCancelled
	}

	tx := gc.DB.WithContext(ctx).Create(&commit)

	if tx.Error != nil {
		if strings.Contains(tx.Error.Error(), `duplicate key value violates unique constraint "idx_commits_commit_id"`) {
			log.Warn().Msgf("already saved commit-id:%s", commit.CommitID)
			return &commit, nil
		} else {
			log.Info().Msgf("getting the save commit error and returning it")
		}
		return &commit, tx.Error
	}
	return &commit, nil
}

// GetAllCommitsByRepositoryName fetches all stores commits by repository name
func (gc *PostgresGitCommitRepository) AllCommitsByRepository(ctx context.Context, repo domain.RepoMetadata, query dtos.APIPagingDto) (*dtos.AllCommitsResponse, error) {
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

func (gc *PostgresGitCommitRepository) TopCommitAuthorsByRepository(ctx context.Context, repo domain.RepoMetadata, limit int) ([]dtos.AuthorCommitCount, error) {
	var results []dtos.AuthorCommitCount
	err := gc.DB.WithContext(ctx).Model(&domain.Commit{}).
		Select("author, COUNT(author) as commit_count").
		Where("repository_name = ?", repo.Name).
		Group("author").
		Order("commit_count DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
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
