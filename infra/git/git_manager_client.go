package git

import (
	"context"
	"time"

	"github.com/kenmobility/git-api-service/internal/domains"
)

type GitManagerClient interface {
	FetchRepoMetadata(ctx context.Context, repositoryName string) (*domains.RepoMetadata, error)
	FetchCommits(ctx context.Context, repo domains.RepoMetadata, since time.Time, until time.Time, lastFetchedCommit string, page, perPage int) ([]domains.Commit, bool, error)
}
