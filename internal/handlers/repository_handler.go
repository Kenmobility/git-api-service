package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kenmobility/github-api-service/common/helpers"
	"github.com/kenmobility/github-api-service/common/message"
	"github.com/kenmobility/github-api-service/common/response"
	"github.com/kenmobility/github-api-service/internal/dtos"
	"github.com/kenmobility/github-api-service/internal/usecases"
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

	repo, err := rh.gitRepositoryUsecase.AddRepository(ctx, input)
	if err != nil {
		response.Failure(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.Success(ctx, http.StatusCreated, "New Repository added successfully", repo)
}

func (rh RepositoryHandlers) FetchAllRepositories(ctx *gin.Context) {
	repos, err := rh.gitRepositoryUsecase.GellAllRepositories(ctx)
	if err != nil {
		response.Failure(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.Success(ctx, http.StatusOK, "successfully fetched all repos", repos)
}
