package persistence

import (
	"context"

	"github.com/kenmobility/github-api-service/common/message"
	"github.com/kenmobility/github-api-service/internal/domains/models"
	"github.com/kenmobility/github-api-service/internal/domains/services"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type GormRepoMetadataRepository struct {
	DB *gorm.DB
}

func NewGormRepoMetadataRepository(db *gorm.DB) services.RepoMetadataRepository {
	return &GormRepoMetadataRepository{DB: db}
}

func (r *GormRepoMetadataRepository) SaveRepoMetadata(ctx context.Context, repo models.RepoMetadata) (*models.RepoMetadata, error) {

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

func (r *GormRepoMetadataRepository) RepoMetadataByPublicId(ctx context.Context, publicId string) (*models.RepoMetadata, error) {
	var repo Repository
	err := r.DB.WithContext(ctx).Where("public_id = ?", publicId).Find(&repo).Error

	if repo.ID == 0 {
		return nil, message.ErrNoRecordFound
	}
	return repo.ToDomain(), err
}

func (r *GormRepoMetadataRepository) RepoMetadataByName(ctx context.Context, name string) (*models.RepoMetadata, error) {
	var repo Repository
	err := r.DB.WithContext(ctx).Where("name = ?", name).Find(&repo).Error
	if repo.ID == 0 {
		return nil, message.ErrNoRecordFound
	}
	return repo.ToDomain(), err
}

func (r *GormRepoMetadataRepository) AllRepoMetadata(ctx context.Context) ([]models.RepoMetadata, error) {
	var dbRepositories []Repository

	err := r.DB.WithContext(ctx).Find(&dbRepositories).Error

	if err != nil {
		return nil, err
	}

	var repoMetaDataResponse []models.RepoMetadata

	for _, dbRepository := range dbRepositories {
		repoMetaDataResponse = append(repoMetaDataResponse, *dbRepository.ToDomain())
	}
	return repoMetaDataResponse, err
}

func (r *GormRepoMetadataRepository) UpdateRepoMetadata(ctx context.Context, repo models.RepoMetadata) (*models.RepoMetadata, error) {
	dbRepo := FromDomainRepo(&repo)

	err := r.DB.WithContext(ctx).Model(&Repository{}).Where(&Repository{PublicID: repo.PublicID}).Updates(&dbRepo).Error
	if err != nil {
		log.Err(err).Msgf("Persistence::UpdateRepoMetadaa error: %v, (%v)", err.Error(), err.Error())
		return nil, err
	}

	return dbRepo.ToDomain(), nil
}
