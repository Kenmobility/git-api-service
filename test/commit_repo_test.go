package test

import (
	"context"
	"testing"

	"github.com/kenmobility/git-api-service/internal/domain"
	"github.com/kenmobility/git-api-service/mocks"
	"github.com/kenmobility/git-api-service/pkg/helpers"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSaveCommitRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mocks.NewMockStore(ctrl)

	commitData := randomCommitdata()

	store.EXPECT().
		SaveCommit(gomock.Any(), commitData).
		Times(1).
		Return(&commitData, nil)

	sCommit, err := store.SaveCommit(context.Background(), commitData)

	require.NoError(t, err)
	require.Equal(t, "sample/repo", sCommit.RepositoryName)
}

func randomCommitdata() domain.Commit {

	repoName := "sample/repo"

	return domain.Commit{
		CommitID:       helpers.RandomString(20),
		Message:        helpers.RandomWords(10),
		URL:            helpers.RandomRepositoryUrl(),
		Author:         helpers.RandomWords(2),
		RepositoryName: repoName,
	}
}
