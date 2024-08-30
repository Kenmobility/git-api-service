package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kenmobility/github-api-service/internal/handlers"
)

func CommitRoutes(r *gin.Engine, ch *handlers.CommitHandlers) {
	r.GET("/commits/:repoId", ch.GetCommitsByRepositoryId)
	r.GET("/top-authors/:repoId", ch.GetTopCommitAuthors)
}
