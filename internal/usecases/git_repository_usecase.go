package usecases

import (
	"context"
	"strings"
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
	StartIndexing(ctx context.Context, input dtos.AddRepositoryRequestDto) (*dtos.GitRepoMetadataResponseDto, error)
	GetById(ctx context.Context, repoId string) (*dtos.GitRepoMetadataResponseDto, error)
	GellAll(ctx context.Context) ([]dtos.GitRepoMetadataResponseDto, error)
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

func (uc *gitRepoUsecase) GetById(ctx context.Context, repoId string) (*dtos.GitRepoMetadataResponseDto, error) {
	repo, err := uc.repoMetadataRepository.RepoMetadataByPublicId(ctx, repoId)
	if err != nil {
		return nil, err
	}

	repoDto := repo.ToDto()

	return &repoDto, nil
}

func (uc *gitRepoUsecase) GellAll(ctx context.Context) ([]dtos.GitRepoMetadataResponseDto, error) {
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

func (uc *gitRepoUsecase) StartIndexing(ctx context.Context, input dtos.AddRepositoryRequestDto) (*dtos.GitRepoMetadataResponseDto, error) {
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

	// Start fetching commits for the new added repository in a goroutine
	go uc.startFetchingAndSavingCommits(ctx, *sRepoMetadata)

	repoDto := sRepoMetadata.ToDto()

	return &repoDto, nil
}

func (uc *gitRepoUsecase) startFetchingAndSavingCommits(ctx context.Context, repo domains.RepoMetadata) {
	page := repo.LastFetchedPage
	lastFetchedCommit := ""
	log.Info().Msgf("fetching commits for repo: %s, starting from page-%d", repo.Name, page)
	for {
		commits, morePages, err := uc.gitClient.FetchCommits(ctx, repo, uc.config.DefaultStartDate, uc.config.DefaultEndDate, "", int(page), uc.config.GitCommitFetchPerPage)
		if err != nil {
			log.Err(err).Msgf("Failed to fetch commits for repository %s: %v", repo.Name, err)
			continue
		}

		// loop through commits and persist each
		for _, commit := range commits {
			_, err := uc.commitRepository.SaveCommit(ctx, commit)
			if err != nil {
				if strings.Contains(err.Error(), `duplicate key value violates unique constraint "idx_commits_commit_id"`) {
					log.Info().Msgf("already saved commit-id:%s for repo %s", commit.CommitID, repo.Name)
					continue
				} else {
					log.Err(err).Msgf("Error saving commit-id:%s for repo %s: %v", commit.CommitID, repo.Name, err)
					continue
				}
			}
			lastFetchedCommit = commit.CommitID
		}

		// Update the repository's last fetched commit in the database
		repo.LastFetchedCommit = lastFetchedCommit
		repo.LastFetchedPage = page
		_, err = uc.repoMetadataRepository.UpdateRepoMetadata(ctx, repo)
		if err != nil {
			log.Debug().Msgf("Error updating repository %s: %v", repo.Name, err)
			continue
		}

		if !morePages {
			break
		}
		page++
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
		go uc.startPeriodicFetching(ctx, repo)
	}
	return nil
}

func (uc *gitRepoUsecase) startPeriodicFetching(ctx context.Context, repo domains.RepoMetadata) error {
	log.Info().Msgf("Commits periodic fetching started for repo %v", repo.Name)
	ticker := time.NewTicker(uc.config.FetchInterval)
	defer ticker.Stop()

	// Initial fetching to start immediately
	uc.fetchAndSaveCommits(ctx, repo)

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Git repo commits fetching service stopped")
			return ctx.Err()
		case <-ticker.C:
			r, err := uc.repoMetadataRepository.RepoMetadataByPublicId(ctx, repo.PublicID)
			if err != nil {
				log.Debug().Msgf("error getting updated last fetched page for repo metadata: %v", err)
				r = &repo
			}
			uc.fetchAndSaveCommits(ctx, *r)
		}
	}
}

func (uc *gitRepoUsecase) fetchAndSaveCommits(ctx context.Context, repo domains.RepoMetadata) {
	log.Info().Msgf("Resume fetching commits for repo: %s", repo.Name)
	page := repo.LastFetchedPage

	repo.IsFetching = true
	uc.repoMetadataRepository.UpdateRepoMetadata(ctx, repo)

	defer func() {
		repo.IsFetching = false
		uc.repoMetadataRepository.UpdateRepoMetadata(ctx, repo)
	}()

	lastFetchedCommit := repo.LastFetchedCommit

	until := uc.config.DefaultEndDate

	for {
		commits, morePages, err := uc.gitClient.FetchCommits(ctx, repo, uc.config.DefaultStartDate, until, lastFetchedCommit, int(page), uc.config.GitCommitFetchPerPage)
		if err != nil {
			log.Error().Msgf("Error fetching commits for repo %s: %v", repo.Name, err)
			return
		}

		if len(commits) == 0 {
			log.Error().Msgf("No new commits for repo %s", repo.Name)
			page = 1               //reset the page
			lastFetchedCommit = "" //don't use sha endpoint
			continue
		}

		for _, commit := range commits {
			c, err := uc.commitRepository.SaveCommit(ctx, commit)
			if err != nil {
				if strings.Contains(err.Error(), `duplicate key value violates unique constraint "idx_commits_commit_id"`) {
					log.Info().Msgf("already saved commit-id:%s for repo %s", commit.CommitID, repo.Name)
					continue
				} else {
					log.Err(err).Msgf("Error saving commit-id:%s for repo %s: %v", commit.CommitID, repo.Name, err)
					continue
				}
			}
			lastFetchedCommit = c.CommitID
		}

		repo.LastFetchedCommit = lastFetchedCommit
		repo.LastFetchedPage = page
		_, err = uc.repoMetadataRepository.UpdateRepoMetadata(ctx, repo)
		if err != nil {
			log.Debug().Msgf("Error updating repository %s: %v", repo.Name, err)
			return
		}

		if !morePages {
			log.Info().Msgf("no more page to fech for repo: %s", repo.Name)
			break
		}

		page++

		until = time.Now()
	}
}
