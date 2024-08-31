package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kenmobility/git-api-service/common/helpers"
	"github.com/kenmobility/git-api-service/common/message"
	"github.com/kenmobility/git-api-service/infra/config"
	"github.com/kenmobility/git-api-service/infra/git"
	"github.com/kenmobility/git-api-service/internal/domains"
	"github.com/kenmobility/git-api-service/internal/http/dtos"
	"github.com/kenmobility/git-api-service/internal/repository"
	"github.com/rs/zerolog/log"
)

type GitRepositoryUsecase interface {
	AddRepository(ctx context.Context, input dtos.AddRepositoryRequestDto) (*dtos.GitRepoMetadataResponseDto, error)
	GetRepositoryById(ctx context.Context, repoId string) (*dtos.GitRepoMetadataResponseDto, error)
	GellAllRepositories(ctx context.Context) ([]dtos.GitRepoMetadataResponseDto, error)
	ResumeFetching(ctx context.Context) error
}

type gitRepoUsecase struct {
	repoMetadataRepository repository.RepoMetadataRepository
	commitRepository       repository.CommitRepository
	gitClient              git.GitManagerClient
	config                 config.Config
}

func NewGitRepositoryUsecase(repoMetadataRepo repository.RepoMetadataRepository, commitRepo repository.CommitRepository,
	gitClient git.GitManagerClient, config config.Config) GitRepositoryUsecase {
	return &gitRepoUsecase{
		repoMetadataRepository: repoMetadataRepo,
		commitRepository:       commitRepo,
		gitClient:              gitClient,
		config:                 config,
	}
}

func (uc *gitRepoUsecase) GetRepositoryById(ctx context.Context, repoId string) (*dtos.GitRepoMetadataResponseDto, error) {
	repo, err := uc.repoMetadataRepository.RepoMetadataByPublicId(ctx, repoId)
	if err != nil {
		return nil, err
	}

	repoDto := repo.ToDto()

	return &repoDto, nil
}

func (uc *gitRepoUsecase) GellAllRepositories(ctx context.Context) ([]dtos.GitRepoMetadataResponseDto, error) {
	repos, err := uc.repoMetadataRepository.AllRepoMetadata(ctx)
	if err != nil {
		return nil, err
	}
	repoDtoResponse := make([]dtos.GitRepoMetadataResponseDto, 0, len(repos))
	for _, repo := range repos {
		repoDtoResponse = append(repoDtoResponse, repo.ToDto())
	}

	return repoDtoResponse, nil
}

func (uc *gitRepoUsecase) AddRepository(ctx context.Context, input dtos.AddRepositoryRequestDto) (*dtos.GitRepoMetadataResponseDto, error) {
	//validate repository name to ensure it has owner and repo name
	if !helpers.IsRepositoryNameValid(input.Name) {
		return nil, message.ErrInvalidRepositoryName
	}

	// ensure repo does not exist on the db
	repo, err := uc.repoMetadataRepository.RepoMetadataByName(ctx, input.Name)
	if err != nil && err != message.ErrNoRecordFound {
		return nil, err
	}

	if repo != nil && repo.Name != "" {
		return nil, message.ErrRepoAlreadyAdded
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
	go uc.startFetchingRepositoryCommits(ctx, *repoMetadata)

	repoDto := sRepoMetadata.ToDto()

	return &repoDto, nil
}

func (uc *gitRepoUsecase) startFetchingRepositoryCommits(ctx context.Context, repo domains.RepoMetadata) {
	ticker := time.NewTicker(uc.config.FetchInterval)
	defer ticker.Stop()

	for range ticker.C {
		// Fetch commits for the repository
		commits, err := uc.gitClient.FetchCommits(ctx, repo, uc.config.DefaultStartDate, uc.config.DefaultEndDate, "")
		if err != nil {
			log.Info().Msgf("Failed to fetch commits for repository %s: %v", repo.Name, err)
			continue
		}

		// loop through commits and persist each
		for _, commit := range commits {
			_, err := uc.commitRepository.SaveCommit(ctx, commit)
			if err != nil {
				log.Info().Msgf("failed to save commitId - %s for repository %s: %v", commit.CommitID, repo.Name, err)
			}
		}
	}
}

func (uc *gitRepoUsecase) ResumeFetching(ctx context.Context) error {
	log.Info().Msg("Resume fetching started ")
	repos, err := uc.repoMetadataRepository.AllRepoMetadata(ctx)
	if err != nil {
		log.Info().Msgf("Error fetching repositories from database: %v", err)
		return err
	}
	log.Info().Msgf("Saved repos %v", repos)

	for _, repo := range repos {
		// Start fetching commits from the last fetched commit
		go log.Info().Err(uc.startPeriodicFetching(ctx, repo))
	}
	return nil
}

func (uc *gitRepoUsecase) startPeriodicFetching(ctx context.Context, repo domains.RepoMetadata) error {
	log.Info().Msgf("Commits periodic fetching started for repo %v", repo.Name)
	ticker := time.NewTicker(uc.config.FetchInterval)
	defer ticker.Stop()

	// Initial fetch to start immediately
	uc.fetchCommits(ctx, repo)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			uc.fetchCommits(ctx, repo)
		}
	}
}

func (uc *gitRepoUsecase) fetchCommits(ctx context.Context, repo domains.RepoMetadata) {
	log.Info().Msgf("fetchcommits for repo: %s", repo.Name)

	repo.IsFetching = true
	uc.repoMetadataRepository.UpdateRepoMetadata(ctx, repo) // Mark as fetching in the DB

	defer func() {
		repo.IsFetching = false
		uc.repoMetadataRepository.UpdateRepoMetadata(ctx, repo) // Mark as not fetching when done
	}()

	lastFetchedCommit := repo.LastFetchedCommit

	var since time.Time
	var until time.Time = time.Now()

	// Fetch commits starting from the last fetched commit
	for {
		commits, err := uc.gitClient.FetchCommits(ctx, repo, since, until, lastFetchedCommit)
		if err != nil {

			log.Info().Msgf("Error fetching commits for repo %s: %v", repo.Name, err)
			return
		}

		if len(commits) == 0 {
			log.Info().Msgf("No new commits for repo %s", repo.Name)
			return
		}

		/*
			for _, commit := range commits {
				if err := uc.commitRepo.SaveCommit(&commit); err != nil {
					log.Printf("Error saving commit for repo %s: %v", repo.Name, err)
					return
				}
				lastFetchedCommit = commit.SHA
			}
		*/

		/* Update the repository's last fetched commit
		repo.LastFetchedCommit = lastFetchedCommit
		if err := uc.repoRepo.UpdateRepository(repo); err != nil {
			log.Printf("Error updating repository %s: %v", repo.Name, err)
			return
		}
		*/
	}
}
