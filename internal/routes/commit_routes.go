package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kenmobility/github-api-service/internal/handlers"
)

func CommitRoutes(r *gin.Engine, ch *handlers.CommitHandlers) {
	r.GET("/repos/:repoId/commits", ch.GetCommitsByRepositoryId)
	r.GET("/repos/:repoId/top-authors", ch.GetTopCommitAuthors)
}
