package services

import (
	"context"

	"github.com/kenmobility/github-api-service/internal/domains/models"
	"github.com/kenmobility/github-api-service/internal/dtos"
)

type CommitRepository interface {
	SaveCommit(ctx context.Context, commit models.Commit) (*models.Commit, error)
	AllCommitsByRepository(ctx context.Context, repoMetadata models.RepoMetadata, query dtos.APIPagingDto) (*dtos.AllCommitsResponse, error)
	TopCommitAuthorsByRepository(ctx context.Context, repo models.RepoMetadata, limit int) ([]string, error)
}
