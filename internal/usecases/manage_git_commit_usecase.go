package usecases

import (
	"context"

	"github.com/kenmobility/git-api-service/internal/http/dtos"
	"github.com/kenmobility/git-api-service/internal/repository"
)

type ManageGitCommitUsecase interface {
	GetAllCommitsByRepository(ctx context.Context, repoId string, query dtos.APIPagingDto) (*string, *dtos.AllCommitsResponse, error)
	GetTopRepositoryCommitAuthors(ctx context.Context, repoId string, limit int) (*string, []dtos.AuthorCommitCount, error)
}

type manageGitCommitUsecase struct {
	commitRepository       repository.CommitRepository
	repoMetadataRepository repository.RepoMetadataRepository
}

func NewManageGitCommitUsecase(commitRepo repository.CommitRepository, repoMetadataRepository repository.RepoMetadataRepository) ManageGitCommitUsecase {
	return &manageGitCommitUsecase{
		commitRepository:       commitRepo,
		repoMetadataRepository: repoMetadataRepository,
	}
}

func (uc *manageGitCommitUsecase) GetAllCommitsByRepository(ctx context.Context, repoId string, query dtos.APIPagingDto) (*string, *dtos.AllCommitsResponse, error) {
	repoMetaData, err := uc.repoMetadataRepository.RepoMetadataByPublicId(ctx, repoId)
	if err != nil {
		return nil, nil, err
	}

	commitsResp, err := uc.commitRepository.AllCommitsByRepository(ctx, *repoMetaData, query)
	if err != nil {
		return nil, nil, err
	}

	return &repoMetaData.Name, commitsResp, nil
}

func (uc *manageGitCommitUsecase) GetTopRepositoryCommitAuthors(ctx context.Context, repoId string, limit int) (*string, []dtos.AuthorCommitCount, error) {
	repoMetaData, err := uc.repoMetadataRepository.RepoMetadataByPublicId(ctx, repoId)
	if err != nil {
		return nil, nil, err
	}

	authors, err := uc.commitRepository.TopCommitAuthorsByRepository(ctx, *repoMetaData, limit)
	if err != nil {
		return nil, nil, err
	}

	return &repoMetaData.Name, authors, nil
}
