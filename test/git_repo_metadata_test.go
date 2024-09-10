package test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kenmobility/git-api-service/common/helpers"
	"github.com/kenmobility/git-api-service/internal/domain"
	"github.com/kenmobility/git-api-service/internal/usecases"
	mockdb "github.com/kenmobility/git-api-service/mock"
	"go.uber.org/mock/gomock"
)

func TestGetRepoMetadataAPI(t *testing.T) {
	repo := randomRepoMetadata()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mockdb.NewMockStore(ctrl)

	testMock := usecases.GitTestRepoUseCase{RepoMetadataRepository: store}

	// build stubs
	store.EXPECT().
		RepoMetadataByPublicId(gomock.Any(), repo.PublicID).
		Times(1).
		Return(repo, nil)

	testMock.GetRepositoryById(context.Background(), repo.PublicID)

}

func randomRepoMetadata() domain.RepoMetadata {
	return domain.RepoMetadata{
		PublicID: uuid.New().String(),
		Name:     helpers.RandomRepositoryName(),
		URL:      helpers.RandomRepositoryUrl(),
		Language: "C++",
	}
}
