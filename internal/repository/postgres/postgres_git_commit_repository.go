package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/kenmobility/git-api-service/internal/domain"
	"github.com/kenmobility/git-api-service/internal/repository"
	"github.com/kenmobility/git-api-service/pkg/message"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type PostgresGitCommitRepository struct {
	DB *gorm.DB
}

func NewPostgresGitCommitRepository(db *gorm.DB) repository.CommitRepository {
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

	dbCommit := FromDomainCommit(&commit)

	tx := gc.DB.WithContext(ctx).Create(&dbCommit)

	if tx.Error != nil {
		if strings.Contains(tx.Error.Error(), `duplicate key value violates unique constraint "idx_commits_commit_id"`) {
			log.Warn().Msgf("already saved commit-id:%s", commit.CommitID)
			return nil, tx.Error
		} else {
			log.Info().Msgf("getting the save commit error and returning it")
		}
		return nil, tx.Error
	}
	return dbCommit.ToDomain(), nil
}

// AllCommitsByRepository fetches all stores commits by repository name
func (gc *PostgresGitCommitRepository) AllCommitsByRepository(ctx context.Context, r domain.RepoMetadata, query domain.APIPagingData) ([]domain.Commit, *domain.PagingInfo, error) {
	var dbCommits []Commit

	var count, queryCount int64

	queryInfo, offset := repository.GetQueryPaginationData(query)

	db := gc.DB.WithContext(ctx).Model(&Commit{}).Where(&Commit{RepositoryName: r.Name})

	db.Count(&count)

	db = db.Offset(offset).Limit(queryInfo.Limit).
		Order(fmt.Sprintf("commits.%s %s", queryInfo.Sort, queryInfo.Direction)).
		Find(&dbCommits)
	db.Count(&queryCount)

	if db.Error != nil {
		log.Info().Msgf("fetch commits error %v", db.Error.Error())

		return nil, nil, db.Error
	}

	pagingInfo := repository.PagingInfo(queryInfo, int(count))
	pagingInfo.Count = len(dbCommits)

	return domainCommits(dbCommits), &pagingInfo, nil
}

func (gc *PostgresGitCommitRepository) TopCommitAuthorsByRepository(ctx context.Context, repo domain.RepoMetadata, limit int) ([]domain.AuthorCommitCount, error) {
	var results []domain.AuthorCommitCount
	err := gc.DB.WithContext(ctx).Model(&domain.Commit{}).
		Select("author, COUNT(author) as commit_count").
		Where("repository_name = ?", repo.Name).
		Group("author").
		Order("commit_count DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

func domainCommits(dbCommits []Commit) []domain.Commit {
	if len(dbCommits) == 0 {
		return []domain.Commit{}
	}

	domainCommits := make([]domain.Commit, 0, len(dbCommits))

	for _, c := range dbCommits {
		cr := domain.Commit{
			CommitID:       c.CommitID,
			Message:        c.Message,
			Author:         c.Author,
			Date:           c.Date,
			URL:            c.URL,
			RepositoryName: c.RepositoryName,
			CreatedAt:      c.CreatedAt,
			UpdatedAt:      c.UpdatedAt,
		}

		domainCommits = append(domainCommits, cr)
	}

	return domainCommits
}
