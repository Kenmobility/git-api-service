package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kenmobility/git-api-service/internal/http/handlers"
)

func RepositoryRoutes(r *gin.Engine, rh *handlers.RepositoryHandlers) {
	r.POST("/repository", rh.AddRepository)
	r.GET("/repositories", rh.FetchAllRepositories)
	r.GET("/repository/:repoId", rh.FetchRepository)
}
