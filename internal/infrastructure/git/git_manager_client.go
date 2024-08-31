package git

import (
	"context"
	"time"

	"github.com/kenmobility/github-api-service/internal/domains/models"
)

type GitManagerClient interface {
	FetchRepoMetadata(ctx context.Context, repositoryName string) (*models.RepoMetadata, error)
	FetchCommits(ctx context.Context, repo models.RepoMetadata, since time.Time, until time.Time, lastFetchedCommit string) ([]models.Commit, error)
}
