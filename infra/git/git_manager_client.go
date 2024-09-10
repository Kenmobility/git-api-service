package git

import (
	"context"
	"time"

	"github.com/kenmobility/git-api-service/internal/domain"
)

type GitManagerClient interface {
	FetchRepoMetadata(ctx context.Context, repositoryName string) (*domain.RepoMetadata, error)
	FetchCommits(ctx context.Context, repo domain.RepoMetadata, since time.Time, until time.Time, lastFetchedCommit string, page, perPage int) ([]domain.Commit, bool, error)
}
