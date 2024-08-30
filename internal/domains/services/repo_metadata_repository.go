package services

import (
	"context"

	"github.com/kenmobility/github-api-service/internal/domains/models"
)

type RepoMetadataRepository interface {
	SaveRepoMetadata(ctx context.Context, repository models.RepoMetadata) (*models.RepoMetadata, error)
	RepoMetadataByPublicId(ctx context.Context, publicId string) (*models.RepoMetadata, error)
	RepoMetadataByName(ctx context.Context, name string) (*models.RepoMetadata, error)
	AllRepoMetadata(ctx context.Context) ([]models.RepoMetadata, error)
}
