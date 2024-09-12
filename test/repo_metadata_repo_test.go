package test

import (
	"context"
	"testing"

	"github.com/kenmobility/git-api-service/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetRepoMetadataRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mocks.NewMockStore(ctrl)

	// Define test data
	repoMetadata := randomRepoMetadata()

	// build stubs
	store.EXPECT().
		RepoMetadataByPublicId(gomock.Any(), repoMetadata.PublicID).
		Return(&repoMetadata, nil).
		Times(1)

	repo, err := store.RepoMetadataByPublicId(context.Background(), repoMetadata.PublicID)

	//require results
	require.NoError(t, err)
	require.Equal(t, repoMetadata, *repo)
}

func TestUpdateRepoMetadataRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mocks.NewMockStore(ctrl)

	// Define test data
	repoMetadata := randomRepoMetadata()

	repoMetadata.IsFetching = false

	store.EXPECT().
		UpdateRepoMetadata(gomock.Any(), repoMetadata).
		Times(1).
		Return(&repoMetadata, nil)

	uRepoMetadata, err := store.UpdateRepoMetadata(context.Background(), repoMetadata)

	//require results
	require.NoError(t, err)
	require.Equal(t, false, uRepoMetadata.IsFetching)
}

func TestSaveRepoMetadataRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mocks.NewMockStore(ctrl)

	repoMetadata := randomRepoMetadata()

	store.EXPECT().
		SaveRepoMetadata(gomock.Any(), repoMetadata).
		Times(1).
		Return(&repoMetadata, nil)

	_, err := store.SaveRepoMetadata(context.Background(), repoMetadata)

	require.NoError(t, err)
}
