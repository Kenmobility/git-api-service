package repository

import (
	"context"

	"github.com/kenmobility/git-api-service/common/message"
	"github.com/kenmobility/git-api-service/internal/domains"
	"gorm.io/gorm"
)

type PostgresGitRepoMetadataRepository struct {
	DB *gorm.DB
}

func NewPostgresGitRepoMetadataRepository(db *gorm.DB) RepoMetadataRepository {
	return &PostgresGitRepoMetadataRepository{DB: db}
}

func (r *PostgresGitRepoMetadataRepository) SaveRepoMetadata(ctx context.Context, repo domains.RepoMetadata) (*domains.RepoMetadata, error) {

	dbRepository := Repository{
		PublicID:        repo.PublicID,
		Name:            repo.Name,
		Description:     repo.Description,
		URL:             repo.URL,
		Language:        repo.Language,
		ForksCount:      repo.ForksCount,
		StarsCount:      repo.StarsCount,
		OpenIssuesCount: repo.OpenIssuesCount,
		WatchersCount:   repo.WatchersCount,
		CreatedAt:       repo.CreatedAt,
		UpdatedAt:       repo.UpdatedAt,
	}
	err := r.DB.WithContext(ctx).Create(&dbRepository).Error
	if err != nil {
		return nil, err
	}

	return dbRepository.ToDomain(), err
}

func (r *PostgresGitRepoMetadataRepository) RepoMetadataByPublicId(ctx context.Context, publicId string) (*domains.RepoMetadata, error) {
	var repo Repository
	err := r.DB.WithContext(ctx).Where("public_id = ?", publicId).Find(&repo).Error

	if repo.ID == 0 {
		return nil, message.ErrNoRecordFound
	}
	return repo.ToDomain(), err
}

func (r *PostgresGitRepoMetadataRepository) RepoMetadataByName(ctx context.Context, name string) (*domains.RepoMetadata, error) {
	var repo Repository
	err := r.DB.WithContext(ctx).Where("name = ?", name).Find(&repo).Error
	if repo.ID == 0 {
		return nil, message.ErrNoRecordFound
	}
	return repo.ToDomain(), err
}

func (r *PostgresGitRepoMetadataRepository) AllRepoMetadata(ctx context.Context) ([]domains.RepoMetadata, error) {
	var dbRepositories []Repository

	err := r.DB.WithContext(ctx).Find(&dbRepositories).Error

	if err != nil {
		return nil, err
	}

	var repoMetaDataResponse []domains.RepoMetadata

	for _, dbRepository := range dbRepositories {
		repoMetaDataResponse = append(repoMetaDataResponse, *dbRepository.ToDomain())
	}
	return repoMetaDataResponse, err
}

func (r *PostgresGitRepoMetadataRepository) UpdateRepoMetadata(ctx context.Context, repo domains.RepoMetadata) (*domains.RepoMetadata, error) {
	dbRepo := FromDomainRepo(&repo)

	err := r.DB.WithContext(ctx).Model(&Repository{}).Where(&Repository{PublicID: repo.PublicID}).Updates(&dbRepo).Error
	if err != nil {
		log.Err(err).Msgf("Persistence::UpdateRepoMetadaa error: %v, (%v)", err.Error(), err.Error())
		return nil, err
	}

	return dbRepo.ToDomain(), nil
}
