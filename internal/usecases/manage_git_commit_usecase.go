package usecases

import (
	"context"

	"github.com/kenmobility/github-api-service/internal/domains/services"
	"github.com/kenmobility/github-api-service/internal/dtos"
)

type ManageGitCommitUsecase interface {
	GetAllCommitsByRepository(ctx context.Context, repoId string, query dtos.APIPagingDto) (*string, *dtos.AllCommitsResponse, error)
	GetTopRepositoryCommitAuthors(ctx context.Context, repoId string, limit int) (*string, []string, error)
}

type manageGitCommitUsecase struct {
	commitRepository       services.CommitRepository
	repoMetadataRepository services.RepoMetadataRepository
}

func NewManageGitCommitUsecase(commitRepo services.CommitRepository, repoMetadataRepository services.RepoMetadataRepository) ManageGitCommitUsecase {
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

func (uc *manageGitCommitUsecase) GetTopRepositoryCommitAuthors(ctx context.Context, repoId string, limit int) (*string, []string, error) {
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
