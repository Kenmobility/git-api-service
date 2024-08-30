package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kenmobility/github-api-service/common/response"
	"github.com/kenmobility/github-api-service/internal/usecases"
)

type CommitHandlers struct {
	manageGitCommitUsecase usecases.ManageGitCommitUsecase
}

func NewCommitHandler(manageGitCommitUsecase usecases.ManageGitCommitUsecase) *CommitHandlers {
	return &CommitHandlers{
		manageGitCommitUsecase: manageGitCommitUsecase,
	}
}

func (ch CommitHandlers) GetCommitsByRepositoryId(ctx *gin.Context) {
	query := getPagingInfo(ctx)

	repositoryId := ctx.Param("repoId")

	if repositoryId == "" {
		response.Failure(ctx, http.StatusBadRequest, "repoId is required", nil)
		return
	}

	repoName, commits, err := ch.manageGitCommitUsecase.GetAllCommitsByRepository(ctx, repositoryId, query)

	if err != nil {
		response.Failure(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	msg := fmt.Sprintf("%s commits fetched successfully", *repoName)

	response.Success(ctx, http.StatusOK, msg, commits)
}

func (ch CommitHandlers) GetTopCommitAuthors(ctx *gin.Context) {
	repositoryId := ctx.Param("repoId")

	if repositoryId == "" {
		response.Failure(ctx, http.StatusBadRequest, "repoId is required", nil)
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))

	if err != nil {
		response.Failure(ctx, http.StatusBadRequest, "invalid limit", err)
		return
	}

	repoName, authors, err := ch.manageGitCommitUsecase.GetTopRepositoryCommitAuthors(ctx, repositoryId, limit)
	if err != nil {
		response.Failure(ctx, http.StatusInternalServerError, "error fetching top authors", err)
		return
	}

	msg := fmt.Sprintf("%v top commit authors of %s repository fetched successfully", limit, *repoName)

	response.Success(ctx, http.StatusOK, msg, authors)
}
