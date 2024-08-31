package repository

import (
	"context"

	"github.com/kenmobility/git-api-service/internal/domains"
)

type RepoMetadataRepository interface {
	SaveRepoMetadata(ctx context.Context, repository domains.RepoMetadata) (*domains.RepoMetadata, error)
	UpdateRepoMetadata(ctx context.Context, repo domains.RepoMetadata) (*domains.RepoMetadata, error)
	RepoMetadataByPublicId(ctx context.Context, publicId string) (*domains.RepoMetadata, error)
	RepoMetadataByName(ctx context.Context, name string) (*domains.RepoMetadata, error)
	AllRepoMetadata(ctx context.Context) ([]domains.RepoMetadata, error)
}
