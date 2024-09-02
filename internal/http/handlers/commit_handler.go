package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kenmobility/git-api-service/common/message"
	"github.com/kenmobility/git-api-service/common/response"
	"github.com/kenmobility/git-api-service/internal/usecases"
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

	repoName, resp, err := ch.manageGitCommitUsecase.GetAllCommitsByRepository(ctx, repositoryId, query)
	if err != nil {
		response.Failure(ctx, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	if len(resp.Commits) == 0 {
		msg := fmt.Sprintf("No commits fetched yet for %s repository...", *repoName)
		response.Success(ctx, http.StatusOK, msg, resp)
		return
	}

	msg := fmt.Sprintf("%s repository commits fetched successfully", *repoName)

	response.Success(ctx, http.StatusOK, msg, resp)
}

func (ch CommitHandlers) GetTopCommitAuthors(ctx *gin.Context) {
	repositoryId := ctx.Param("repoId")

	if repositoryId == "" {
		response.Failure(ctx, http.StatusBadRequest, "repoId is required", nil)
		return
	}
	repoName, authors, err := ch.manageGitCommitUsecase.GetTopRepositoryCommitAuthors(ctx, repositoryId, getPagingInfo(ctx).Limit)
	if err != nil {
		if err == message.ErrNoRecordFound {
			response.Failure(ctx, http.StatusBadRequest, message.ErrInvalidRepositoryId.Error(), message.ErrInvalidRepositoryId.Error())
			return
		}
		response.Failure(ctx, http.StatusInternalServerError, "error fetching repo authors", err.Error())
		return
	}

	if len(authors) == 0 {
		response.Success(ctx, http.StatusOK, "no top authors fetched yet", nil)
		return
	}

	msg := fmt.Sprintf("%v top commit authors of %s repository fetched successfully", len(authors), *repoName)

	response.Success(ctx, http.StatusOK, msg, authors)
}
