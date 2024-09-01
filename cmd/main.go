package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kenmobility/git-api-service/common/message"
	"github.com/kenmobility/git-api-service/infra/config"
	"github.com/kenmobility/git-api-service/infra/database"
	"github.com/kenmobility/git-api-service/infra/git"
	"github.com/kenmobility/git-api-service/infra/logger"
	"github.com/kenmobility/git-api-service/internal/http/dtos"
	"github.com/kenmobility/git-api-service/internal/http/handlers"
	"github.com/kenmobility/git-api-service/internal/http/routes"
	"github.com/kenmobility/git-api-service/internal/repository"
	"github.com/kenmobility/git-api-service/internal/usecases"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configures system wide Logger object
	log := logger.New("main")

	// load env variables
	config, err := config.LoadConfig("")
	if err != nil {
		log.Fatal().Msgf("failed to load config %v, (%v)", err.Error(), err.Error())
	}

	// Initialize Database Client
	dbClient := database.NewPostgresDatabase(*config)

	// establish database connection
	db, err := dbClient.ConnectDb()
	if err != nil {
		log.Fatal().Msgf("failed to establish postgres database connection: %v, (%v)", err.Error(), err.Error())
	}

	// Run migrations
	if err := dbClient.Migrate(db); err != nil {
		log.Fatal().Msgf("failed to run database migrations: %v, (%v)", err.Error(), err.Error())
	}

	// Initialize repositories
	commitRepository := repository.NewPostgresGitCommitRepository(db)
	repoMetadataRepository := repository.NewPostgresGitRepoMetadataRepository(db)

	// Initialize Git Manager Client
	gitClient := git.NewGitHubClient(config.GitHubApiBaseURL, config.GitHubToken, config.FetchInterval)

	// Initialize use cases and handlers
	gitCommitUsecase := usecases.NewManageGitCommitUsecase(commitRepository, repoMetadataRepository)
	gitRepositoryUsecase := usecases.NewGitRepositoryUsecase(repoMetadataRepository, commitRepository, gitClient, *config)

	// seed and set 'chromium/chromium' repo as default repository if not seeded
	err = seedDefaultRepository(config, gitRepositoryUsecase)
	if err != nil && err != message.ErrRepoAlreadyAdded {
		log.Fatal().Msgf("failed to seed default repository: %v, (%v)", err.Error(), err.Error())
	}

	// Initialize handlers
	commitHandler := handlers.NewCommitHandler(gitCommitUsecase)
	repositoryHandler := handlers.NewRepositoryHandler(gitRepositoryUsecase)

	ginEngine := gin.Default()

	// register routes
	routes.CommitRoutes(ginEngine, commitHandler)
	routes.RepositoryRoutes(ginEngine, repositoryHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Address, config.Port),
		Handler: ginEngine,
	}

	// create a context with cancellation to gracefully shut down Git commits monitoring if server shuts down
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// Resume repo commits fetching for all repositories
	go gitRepositoryUsecase.ResumeFetching(ctx)

	// start web server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Msgf("listen: %v, (%v)\n", err.Error(), err.Error())
	}
	log.Info().Msgf("Git API Service is listening on address %s", server.Addr)
}

// SeedRepository seeds a default chromium repo and set it as tracking
func seedDefaultRepository(config *config.Config, repositoryUsecase usecases.GitRepositoryUsecase) error {
	defaultRepo := dtos.AddRepositoryRequestDto{
		Name: config.DefaultRepository,
	}
	repo, err := repositoryUsecase.AddRepository(context.Background(), defaultRepo)
	if err != nil && err != message.ErrNoRecordFound {
		return err
	}

	if repo != nil {
		log.Info().Msgf("Successfully seeded default repository: %s", repo.Name)
	}
	return err
}
