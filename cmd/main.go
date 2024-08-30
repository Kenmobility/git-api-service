package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kenmobility/github-api-service/config"
	"github.com/kenmobility/github-api-service/database"
	"github.com/kenmobility/github-api-service/internal/domains/models"
	"github.com/kenmobility/github-api-service/internal/domains/services"
	"github.com/kenmobility/github-api-service/internal/handlers"
	"github.com/kenmobility/github-api-service/internal/infrastructure/git"
	"github.com/kenmobility/github-api-service/internal/infrastructure/persistence"
	"github.com/kenmobility/github-api-service/internal/routes"
	"github.com/kenmobility/github-api-service/internal/usecases"
)

func main() {
	// load env variables
	config, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// establish database connection
	dbClient := database.NewPostgresDatabase(*config)

	db, err := dbClient.ConnectDb()
	if err != nil {
		log.Fatalf("failed to establish postgres database connection: %v", err)
	}

	// Run migrations
	if err := dbClient.Migrate(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Initialize repositories
	commitRepo := persistence.NewGormCommitRepository(db)
	repoMetadataRepo := persistence.NewGormRepositoRepository(db)

	// seed and set 'chromium/chromium' repo as default repository to track
	err = seedDefaultRepository(config, repoMetadataRepo)
	if err != nil {
		log.Fatalf("failed to seed default repository to database: %v", err)
	}

	gitClient := git.NewGitHubClient(config.GitHubApiBaseURL, config.GitHubToken, config.FetchInterval, commitRepo, repoMetadataRepo)

	// Initialize use cases and handlers
	gitCommitUsecase := usecases.NewManageGitCommitUsecase(commitRepo, repoMetadataRepo)
	gitRepositoryUsecase := usecases.NewGitRepositoryUsecase(repoMetadataRepo, commitRepo, gitClient, *config)

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

	// start web server
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v\n", err)
		}
		log.Printf("GitHub Service is listening on address %s", server.Addr)
	}()

	// create a context with cancellation to gracefully shut down GitHub tracker service if server shuts down
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown: err=%v", err)
	}
}

// SeedRepository seeds a default chromium repo and set it as tracking
func seedDefaultRepository(config *config.Config, repoRepo services.RepoMetadataRepository) error {
	defaultRepo := models.RepoMetadata{
		PublicID: uuid.New().String(),
		Name:     config.DefaultRepository,
		URL:      "https://github.com/chromium/chromium",
		Language: "C++",
	}

	_, err := repoRepo.RepoMetadataByName(context.Background(), defaultRepo.Name)
	if err != nil {
		log.Printf("Repository %s already exists in the database, skipping seeding.", defaultRepo.Name)
		return err
	}

	_, err = repoRepo.SaveRepoMetadata(context.Background(), defaultRepo)
	if err != nil {
		log.Printf("failed to see default repository: %v", err)
		return err
	}

	log.Printf("Successfully seeded default repository: %s", defaultRepo.Name)
	return err
}
