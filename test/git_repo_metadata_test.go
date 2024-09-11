package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/kenmobility/git-api-service/common/helpers"
	"github.com/kenmobility/git-api-service/internal/domain"
	mockdb "github.com/kenmobility/git-api-service/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetRepoMetadataAPI(t *testing.T) {
	repo := randomRepoMetadata()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mockdb.NewMockStore(ctrl)

	// build stubs
	store.EXPECT().
		RepoMetadataByPublicId(gomock.Any(), repo.PublicID).
		Times(1).
		Return(repo, nil)

	// Start a local HTTP test server and send request
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()

	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/repository/%s", repo.PublicID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

}

func randomRepoMetadata() domain.RepoMetadata {
	return domain.RepoMetadata{
		PublicID: uuid.New().String(),
		Name:     helpers.RandomRepositoryName(),
		URL:      helpers.RandomRepositoryUrl(),
		Language: "C++",
	}
}
