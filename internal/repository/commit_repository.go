package repository

import (
	"context"

	"github.com/kenmobility/git-api-service/internal/domain"
	"github.com/kenmobility/git-api-service/internal/http/dtos"
)

type CommitRepository interface {
	SaveCommit(ctx context.Context, commit domain.Commit) (*domain.Commit, error)
	AllCommitsByRepository(ctx context.Context, repoMetadata domain.RepoMetadata, query dtos.APIPagingDto) (*dtos.AllCommitsResponse, error)
	TopCommitAuthorsByRepository(ctx context.Context, repo domain.RepoMetadata, limit int) ([]dtos.AuthorCommitCount, error)
}
