package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kenmobility/git-api-service/common/message"
	"github.com/kenmobility/git-api-service/infra/config"
	"github.com/kenmobility/git-api-service/infra/database"
	"github.com/kenmobility/git-api-service/infra/git"
	"github.com/kenmobility/git-api-service/internal/http/dtos"
	"github.com/kenmobility/git-api-service/internal/http/handlers"
	"github.com/kenmobility/git-api-service/internal/http/routes"
	"github.com/kenmobility/git-api-service/internal/repository"
	"github.com/kenmobility/git-api-service/internal/usecases"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configures system wide Logger object
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Caller().Logger()
	// make logger human-readable, only locally
	if os.Getenv("APP_ENV") == "" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
	}

	// load env variables
	config, err := config.LoadConfig("")
	if err != nil {
		log.Fatal().Msgf("failed to load config %v, (%v)", err.Error(), err.Error())
	}

	// establish database connection
	dbClient := database.NewPostgresDatabase(*config)
	db, err := dbClient.ConnectDb()
	if err != nil {
		log.Fatal().Msgf("failed to establish postgres database connection: %v, (%v)", err.Error(), err.Error())
	}

	// Run database migrations
	if err := dbClient.Migrate(db); err != nil {
		log.Fatal().Msgf("failed to run database migrations: %v, (%v)", err.Error(), err.Error())
	}

	// Initialize various layers
	commitRepository := repository.NewPostgresGitCommitRepository(db)
	repoMetadataRepository := repository.NewPostgresGitRepoMetadataRepository(db)

	gitClient := git.NewGitHubClient(config.GitHubApiBaseURL, config.GitHubToken, config.FetchInterval)

	gitCommitUsecase := usecases.NewManageGitCommitUsecase(commitRepository, repoMetadataRepository)
	gitRepositoryUsecase := usecases.NewGitRepositoryUsecase(repoMetadataRepository, commitRepository, gitClient, *config)

	commitHandler := handlers.NewCommitHandler(gitCommitUsecase)
	repositoryHandler := handlers.NewRepositoryHandler(gitRepositoryUsecase)

	//seed default repo
	err = seedDefaultRepository(config, gitRepositoryUsecase)
	if err != nil && err != message.ErrRepoAlreadyAdded {
		log.Fatal().Msgf("failed to seed default repository: %v, (%v)", err.Error(), err.Error())
	}

	ginEngine := gin.Default()

	// register routes
	routes.CommitRoutes(ginEngine, commitHandler)
	routes.RepositoryRoutes(ginEngine, repositoryHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Address, config.Port),
		Handler: ginEngine,
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// Resume repo commits fetching for all saved repositories
	go gitRepositoryUsecase.ResumeFetching(ctx)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Msgf("listen: %v, (%v)\n", err.Error(), err.Error())
	}
	log.Info().Msgf("Git API Service is listening on address %s", server.Addr)
}

// seedDefaultRepository seeds a default repository to database
func seedDefaultRepository(config *config.Config, repositoryUsecase usecases.GitRepositoryUsecase) error {
	defaultRepo := dtos.AddRepositoryRequestDto{
		Name: config.DefaultRepository,
	}
	repo, err := repositoryUsecase.StartIndexing(context.Background(), defaultRepo)
	if err != nil && err != message.ErrNoRecordFound {
		return err
	}

	if repo != nil {
		log.Info().Msgf("Successfully seeded default repository: %s", repo.Name)
	}
	return err
}
