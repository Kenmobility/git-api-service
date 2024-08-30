package usecases

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kenmobility/github-api-service/common/helpers"
	"github.com/kenmobility/github-api-service/common/message"
	"github.com/kenmobility/github-api-service/config"
	"github.com/kenmobility/github-api-service/internal/domains/models"
	"github.com/kenmobility/github-api-service/internal/domains/services"
	"github.com/kenmobility/github-api-service/internal/dtos"
	"github.com/kenmobility/github-api-service/internal/infrastructure/git"
)

type GitRepositoryUsecase interface {
	AddRepository(ctx context.Context, input dtos.AddRepositoryRequestDto) (*models.RepoMetadata, error)
	GetRepositoryById(ctx context.Context, repoId string) (*models.RepoMetadata, error)
	GellAllRepositories(ctx context.Context) ([]models.RepoMetadata, error)
}

type addGitRepoUsecase struct {
	repoMetadataRepository services.RepoMetadataRepository
	commitRepository       services.CommitRepository
	gitClient              git.GitManagerClient
	config                 config.Config
}

func NewGitRepositoryUsecase(repoMetadataRepo services.RepoMetadataRepository, commitRepo services.CommitRepository,
	gitClient git.GitManagerClient, config config.Config) GitRepositoryUsecase {
	return &addGitRepoUsecase{
		repoMetadataRepository: repoMetadataRepo,
		commitRepository:       commitRepo,
		gitClient:              gitClient,
		config:                 config,
	}
}

func (uc *addGitRepoUsecase) GetRepositoryById(ctx context.Context, repoId string) (*models.RepoMetadata, error) {
	return uc.repoMetadataRepository.RepoMetadataByPublicId(ctx, repoId)
}

func (uc *addGitRepoUsecase) GellAllRepositories(ctx context.Context) ([]models.RepoMetadata, error) {
	return uc.repoMetadataRepository.AllRepoMetadata(ctx)
}

func (uc *addGitRepoUsecase) AddRepository(ctx context.Context, input dtos.AddRepositoryRequestDto) (*models.RepoMetadata, error) {
	//validate repository name to ensure it has owner and repo name
	if !helpers.IsRepositoryNameValid(input.Name) {
		return nil, message.ErrInvalidRepositoryName
	}

	// try fetching repo meta data from GitManagerClient to ensure repository with such name exists
	repoMetadata, err := uc.gitClient.FetchRepoMetadata(ctx, input.Name)
	if err != nil {
		return nil, err
	}

	// update other repository metadata
	repoMetadata.PublicID = uuid.New().String()
	repoMetadata.CreatedAt = time.Now()
	repoMetadata.UpdatedAt = time.Now()

	sRepoMetadata, err := uc.repoMetadataRepository.SaveRepoMetadata(ctx, *repoMetadata)
	if err != nil {
		return nil, err
	}

	// Start fetching commits for the new added repository in a new gorouting
	go uc.startFetchingRepositoryCommits(ctx, repoMetadata.Name)

	return sRepoMetadata, nil
}

func (uc *addGitRepoUsecase) startFetchingRepositoryCommits(ctx context.Context, repoName string) {
	ticker := time.NewTicker(uc.config.FetchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Fetch commits for the repository
			commits, err := uc.gitClient.FetchCommits(ctx, repoName, uc.config.DefaultStartDate, uc.config.DefaultEndDate)
			if err != nil {
				log.Printf("Failed to fetch commits for repository %s: %v", repoName, err)
				continue
			}

			// loop through commits and persist each
			for _, commit := range commits {
				_, err := uc.commitRepository.SaveCommit(ctx, commit)
				if err != nil {
					log.Printf("failed to save commitId - %s for repository %s: %v", commit.CommitID, repoName, err)
				}
			}

		}
	}
}
