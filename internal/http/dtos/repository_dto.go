package dtos

import (
	"time"

	"github.com/kenmobility/git-api-service/internal/domain"
)

type AddRepositoryRequestDto struct {
	Name string `json:"name" validate:"required"`
}

type GitRepoMetadataResponseDto struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	URL             string `json:"url"`
	Language        string `json:"language"`
	ForksCount      int    `json:"forks_count"`
	StarsCount      int    `json:"stars_count"`
	OpenIssuesCount int    `json:"open_issues_count"`
	WatchersCount   int    `json:"watchers_count"`
	CreatedAt       string `json:"added_at"`
	UpdatedAt       string `json:"last_updated_at"`
}

func RepoMetadataResponse(r domain.RepoMetadata) GitRepoMetadataResponseDto {
	return GitRepoMetadataResponseDto{
		Id:              r.PublicID,
		Name:            r.Name,
		Description:     r.Description,
		URL:             r.URL,
		Language:        r.Language,
		ForksCount:      r.ForksCount,
		StarsCount:      r.StarsCount,
		OpenIssuesCount: r.OpenIssuesCount,
		WatchersCount:   r.WatchersCount,
		CreatedAt:       r.CreatedAt.Format(time.RFC850),
		UpdatedAt:       r.UpdatedAt.Format(time.RFC850),
	}
}

func AllRepoMetadataResponse(repos []domain.RepoMetadata) []GitRepoMetadataResponseDto {
	if len(repos) == 0 {
		return []GitRepoMetadataResponseDto{}
	}

	reposResponse := make([]GitRepoMetadataResponseDto, 0, len(repos))

	for _, r := range repos {
		rr := GitRepoMetadataResponseDto{
			Id:              r.PublicID,
			Name:            r.Name,
			Description:     r.Description,
			URL:             r.URL,
			Language:        r.Language,
			ForksCount:      r.ForksCount,
			StarsCount:      r.StarsCount,
			OpenIssuesCount: r.OpenIssuesCount,
			WatchersCount:   r.WatchersCount,
			CreatedAt:       r.CreatedAt.Format(time.RFC850),
			UpdatedAt:       r.UpdatedAt.Format(time.RFC850),
		}

		reposResponse = append(reposResponse, rr)
	}

	return reposResponse
}
