package repository

import (
	"context"

	"github.com/kenmobility/git-api-service/internal/domain"
)

type CommitRepository interface {
	SaveCommit(ctx context.Context, commit domain.Commit) (*domain.Commit, error)
	GetByCommitID(ctx context.Context, commitID string) (*domain.Commit, error)
	AllCommitsByRepository(ctx context.Context, repoMetadata domain.RepoMetadata, query domain.APIPaging) ([]domain.Commit, *domain.PagingInfo, error)
	TopCommitAuthorsByRepository(ctx context.Context, repo domain.RepoMetadata, limit int) ([]domain.AuthorCommitCount, error)
}
