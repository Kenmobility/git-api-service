package repository

import (
	"context"

	"github.com/kenmobility/git-api-service/internal/domains"
	"github.com/kenmobility/git-api-service/internal/http/dtos"
)

type CommitRepository interface {
	SaveCommit(ctx context.Context, commit domains.Commit) (*domains.Commit, error)
	AllCommitsByRepository(ctx context.Context, repoMetadata domains.RepoMetadata, query dtos.APIPagingDto) (*dtos.AllCommitsResponse, error)
	TopCommitAuthorsByRepository(ctx context.Context, repo domains.RepoMetadata, limit int) ([]string, error)
}
