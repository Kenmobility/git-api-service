package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kenmobility/git-api-service/infra/config"
	"github.com/kenmobility/git-api-service/infra/git"
	"github.com/kenmobility/git-api-service/internal/domain"
	"github.com/kenmobility/git-api-service/internal/repository"
	"github.com/kenmobility/git-api-service/pkg/helpers"
	"github.com/kenmobility/git-api-service/pkg/message"
	"github.com/rs/zerolog/log"
)

type GitRepositoryUsecase interface {
	StartIndexing(ctx context.Context, repositoryName string) (*domain.RepoMetadata, error)
	GetById(ctx context.Context, repoId string) (*domain.RepoMetadata, error)
	GellAll(ctx context.Context) ([]domain.RepoMetadata, error)
	ResumeFetching(ctx context.Context) error
	UpdateFetchingStatusForAllRepositories(ctx context.Context, status bool) error
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

func (uc *gitRepoUsecase) GetById(ctx context.Context, repoId string) (*domain.RepoMetadata, error) {
	repo, err := uc.repoMetadataRepository.RepoMetadataByPublicId(ctx, repoId)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (uc *gitRepoUsecase) GellAll(ctx context.Context) ([]domain.RepoMetadata, error) {
	repos, err := uc.repoMetadataRepository.AllRepoMetadata(ctx)
	if err != nil {
		return nil, err
	}
	repoDtoResponse := make([]domain.RepoMetadata, 0, len(repos))
	repoDtoResponse = append(repoDtoResponse, repos...)

	return repoDtoResponse, nil
}

func (uc *gitRepoUsecase) StartIndexing(ctx context.Context, repositoryName string) (*domain.RepoMetadata, error) {
	//validate repository name to ensure it has owner and repo name
	if !helpers.IsRepositoryNameValid(repositoryName) {
		return nil, message.ErrInvalidRepositoryName
	}

	// ensure repo does not exist on the db
	repo, err := uc.repoMetadataRepository.RepoMetadataByName(ctx, repositoryName)
	if err != nil && err != message.ErrNoRecordFound {
		return nil, err
	}

	if repo != nil && repo.Name != "" {
		return nil, message.ErrRepoAlreadyAdded
	}

	repoMetadata, err := uc.gitClient.FetchRepoMetadata(ctx, repositoryName)
	if err != nil {
		return nil, err
	}

	// update other repository metadata
	repoMetadata.PublicID = uuid.New().String()
	repoMetadata.CreatedAt = time.Now()
	repoMetadata.UpdatedAt = time.Now()
	repoMetadata.IsFetching = true

	sRepoMetadata, err := uc.repoMetadataRepository.SaveRepoMetadata(ctx, *repoMetadata)
	if err != nil {
		return nil, err
	}

	// Start fetching commits for the new added repository in a Goroutine
	go uc.startRepoIndexing(ctx, *sRepoMetadata)

	return sRepoMetadata, nil
}

func (uc *gitRepoUsecase) startRepoIndexing(ctx context.Context, repo domain.RepoMetadata) {
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
				log.Err(err).Msgf("error saving commit-id:%s for repo %s", commit.CommitID, repo.Name)
				continue
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
			// update isFetching to false as flag for start of monitoring
			repo.IsFetching = false
			_, err = uc.repoMetadataRepository.UpdateRepoMetadata(ctx, repo)
			if err != nil {
				log.Err(err).Msgf("Error updating isFetching column of repository %s: %v", repo.Name, err)
			}
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

func (uc *gitRepoUsecase) startPeriodicFetching(ctx context.Context, repo domain.RepoMetadata) error {
	ticker := time.NewTicker(uc.config.FetchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Warn().Msgf("Git repository [%s] commits monitoring service stopped", repo.Name)
			return ctx.Err()
		case <-ticker.C:
			r, err := uc.repoMetadataRepository.RepoMetadataByPublicId(ctx, repo.PublicID)
			if err != nil {
				log.Debug().Msgf("error getting repo metadata for monitoring: %v", err)
				return err
			}
			if !r.IsFetching {
				log.Info().Msgf("Commits periodic fetching started for repo %v", repo.Name)
				uc.fetchAndReconcileCommits(ctx, *r)
			}
		}
	}
}

func (uc *gitRepoUsecase) fetchAndReconcileCommits(ctx context.Context, repo domain.RepoMetadata) {
	log.Info().Msgf("Resume fetching and reconciling commits for repo: %s", repo.Name)
	page := repo.LastFetchedPage

	lastFetchedCommit := repo.LastFetchedCommit

	until := uc.config.DefaultEndDate

	for {
		select {
		case <-ctx.Done():
			log.Warn().Msgf("Git repository [%s] fetchAndReconcileCommits service stopped", repo.Name)
			return
		default:
			commits, morePages, err := uc.gitClient.FetchCommits(ctx, repo, uc.config.DefaultStartDate, until, lastFetchedCommit, int(page), uc.config.GitCommitFetchPerPage)
			if err != nil {
				log.Error().Msgf("Error fetching commits for repo %s: %v", repo.Name, err)
				return
			}

			if len(commits) == 0 {
				log.Warn().Msgf("No new commits for repo %s, resetting page to 1", repo.Name)
				page = 1               //reset the page
				lastFetchedCommit = "" //don't use sha endpoint
				continue
			}

			for _, commit := range commits {
				_, err = uc.commitRepository.GetByCommitID(ctx, commit.CommitID)
				if err != nil && err != message.ErrNoRecordFound && err != message.ErrContextCancelled {
					log.Err(err).Msgf("error getting commit by commit-id:%s", commit.CommitID)
				}
				if err == message.ErrNoRecordFound {
					_, err := uc.commitRepository.SaveCommit(ctx, commit)
					if err != nil {
						log.Err(err).Msgf("error saving commit-id:%s for repo %s", commit.CommitID, repo.Name)
						continue
					}
					lastFetchedCommit = commit.CommitID
				}
			}

			repo.LastFetchedCommit = lastFetchedCommit
			repo.LastFetchedPage = page
			_, err = uc.repoMetadataRepository.UpdateRepoMetadata(ctx, repo)
			if err != nil && err != message.ErrContextCancelled {
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
}

func (uc *gitRepoUsecase) UpdateFetchingStatusForAllRepositories(ctx context.Context, status bool) error {
	return uc.repoMetadataRepository.UpdateFetchingStateForAllRepos(ctx, status)
}
