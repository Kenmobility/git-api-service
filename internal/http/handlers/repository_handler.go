package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kenmobility/git-api-service/internal/http/dtos"
	"github.com/kenmobility/git-api-service/internal/usecases"
	"github.com/kenmobility/git-api-service/pkg/helpers"
	"github.com/kenmobility/git-api-service/pkg/message"
	"github.com/kenmobility/git-api-service/pkg/response"
)

type RepositoryHandlers struct {
	gitRepositoryUsecase usecases.GitRepositoryUsecase
}

func NewRepositoryHandler(gitRepositoryUsecase usecases.GitRepositoryUsecase) *RepositoryHandlers {
	return &RepositoryHandlers{
		gitRepositoryUsecase: gitRepositoryUsecase,
	}
}

func (rh RepositoryHandlers) AddRepository(ctx *gin.Context) {
	var input dtos.AddRepositoryRequestDto

	err := ctx.BindJSON(&input)
	if err != nil {
		response.Failure(ctx, http.StatusBadRequest, "invalid input", err)
		return
	}

	inputErrors := helpers.ValidateInput(input)
	if inputErrors != nil {
		response.Failure(ctx, http.StatusBadRequest, message.ErrInvalidInput.Error(), inputErrors)
		return
	}

	repo, err := rh.gitRepositoryUsecase.StartIndexing(ctx, input)
	if err != nil {
		if err == message.ErrRepoAlreadyAdded {
			response.Failure(ctx, http.StatusBadRequest, err.Error(), err.Error())
			return
		}
		response.Failure(ctx, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	response.Success(ctx, http.StatusCreated, "Repository successfully indexed, its commits are being fetched...", repo)
}

func (rh RepositoryHandlers) FetchAllRepositories(ctx *gin.Context) {
	repos, err := rh.gitRepositoryUsecase.GellAll(ctx)
	if err != nil {
		response.Failure(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	if len(repos) == 0 {
		response.Success(ctx, http.StatusOK, "no repository indexed yet", repos)
		return
	}
	response.Success(ctx, http.StatusOK, "successfully fetched all repos", repos)
}

func (rh RepositoryHandlers) FetchRepository(ctx *gin.Context) {
	repositoryId := ctx.Param("repoId")
	if repositoryId == "" {
		response.Failure(ctx, http.StatusBadRequest, "repoId is required", nil)
		return
	}

	repo, err := rh.gitRepositoryUsecase.GetById(ctx, repositoryId)
	if err != nil {
		if err == message.ErrNoRecordFound {
			response.Failure(ctx, http.StatusBadRequest, message.ErrInvalidRepositoryId.Error(), message.ErrInvalidRepositoryId.Error())
			return
		}
		response.Failure(ctx, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "successfully fetched repository", repo)
}
