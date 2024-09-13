package usecases

import (
	"context"

	"github.com/kenmobility/git-api-service/internal/domain"
	//"github.com/kenmobility/git-api-service/internal/http/dtos"
	"github.com/kenmobility/git-api-service/internal/repository"
)

type ManageGitCommitUsecase interface {
	GetAllCommitsByRepository(ctx context.Context, repoId string, query domain.APIPaging) (*string, []domain.Commit, *domain.PagingInfo, error)
	GetTopRepositoryCommitAuthors(ctx context.Context, repoId string, limit int) (*string, []domain.AuthorCommitCount, error)
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

func (uc *manageGitCommitUsecase) GetAllCommitsByRepository(ctx context.Context, repoId string, query domain.APIPaging) (*string, []domain.Commit, *domain.PagingInfo, error) {
	repoMetaData, err := uc.repoMetadataRepository.RepoMetadataByPublicId(ctx, repoId)
	if err != nil {
		return nil, nil, nil, err
	}

	commits, pagingInfo, err := uc.commitRepository.AllCommitsByRepository(ctx, *repoMetaData, query)
	if err != nil {
		return nil, nil, nil, err
	}

	return &repoMetaData.Name, commits, pagingInfo, nil
}

func (uc *manageGitCommitUsecase) GetTopRepositoryCommitAuthors(ctx context.Context, repoId string, limit int) (*string, []domain.AuthorCommitCount, error) {
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
