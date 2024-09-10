package repository

import (
	"context"

	"github.com/kenmobility/git-api-service/internal/domain"
)

type RepoMetadataRepository interface {
	SaveRepoMetadata(ctx context.Context, repository domain.RepoMetadata) (*domain.RepoMetadata, error)
	UpdateRepoMetadata(ctx context.Context, repo domain.RepoMetadata) (*domain.RepoMetadata, error)
	RepoMetadataByPublicId(ctx context.Context, publicId string) (*domain.RepoMetadata, error)
	RepoMetadataByName(ctx context.Context, name string) (*domain.RepoMetadata, error)
	AllRepoMetadata(ctx context.Context) ([]domain.RepoMetadata, error)
}
