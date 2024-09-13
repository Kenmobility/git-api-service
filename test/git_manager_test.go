package test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kenmobility/git-api-service/internal/domain"
	"github.com/kenmobility/git-api-service/mocks"
	"github.com/kenmobility/git-api-service/pkg/helpers"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestFetchRepoMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGitClient := mocks.NewMockGitManagerClient(ctrl)

	// Define test data
	repoName := "sample/repo"
	expectedMetadata := domain.RepoMetadata{
		Name:        "sample/repo",
		Description: "A sample repository",
		URL:         "https://example.com/sample/repo",
		StarsCount:  100,
		ForksCount:  200,
	}

	mockGitClient.EXPECT().
		FetchRepoMetadata(gomock.Any(), repoName).
		Return(&expectedMetadata, nil).
		Times(1)

	metadata, err := mockGitClient.FetchRepoMetadata(context.Background(), repoName)

	require.NoError(t, err)
	require.Equal(t, expectedMetadata, *metadata)
}

func TestFetchCommits(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock for the GitManagerClient
	mockGitClient := mocks.NewMockGitManagerClient(ctrl)

	// Define test data
	repoMetadata := randomRepoMetadata()
	since := time.Now().Add(-24 * time.Hour)
	until := time.Now()
	page := 1
	perPage := 2
	lastFetchedCommit := "abc123"
	expectedCommits := []domain.Commit{
		{CommitID: "abc123", Author: "john", Message: "Initial commit", URL: "https://example.com/commit/abc123"},
		{CommitID: "def456", Author: "jane", Message: "Update README", URL: "https://example.com/commit/def456"},
	}

	// Set up the expected behavior for the mock
	mockGitClient.EXPECT().
		FetchCommits(gomock.Any(), repoMetadata, since, until, lastFetchedCommit, page, perPage).
		Return(expectedCommits, false, nil).
		Times(1)

	// Now call the method using the mock
	commits, hasMore, err := mockGitClient.FetchCommits(context.Background(), repoMetadata, since, until, lastFetchedCommit, page, perPage)

	// require results
	require.NoError(t, err)
	require.Equal(t, hasMore, false)
	require.Equal(t, expectedCommits, commits)
}

func randomRepoMetadata() domain.RepoMetadata {
	return domain.RepoMetadata{
		PublicID:   uuid.New().String(),
		Name:       helpers.RandomRepositoryName(),
		URL:        helpers.RandomRepositoryUrl(),
		Language:   "C++",
		IsFetching: true,
	}
}
