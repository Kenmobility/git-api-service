package dtos

import (
	"time"

	"github.com/kenmobility/git-api-service/internal/domain"
)

type AllCommitResponse struct {
	Commits  []CommitResponseDto `json:"commits"`
	PageInfo PagingInfoDto       `json:"page_info"`
}

type CommitResponseDto struct {
	CommitID   string    `json:"commit_id"`
	Message    string    `json:"message"`
	Author     string    `json:"author"`
	Date       time.Time `json:"date"`
	URL        string    `json:"url"`
	Repository string    `json:"repository"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// AuthorCommitCountDto holds the result with author and count of commits
type AuthorCommitCountDto struct {
	Author      string `json:"author"`
	CommitCount int    `json:"commit_count"`
}

// AuthorCommitCountResponse maps to dto response from AuthorCommitCount domain object
func AuthorCommitCountResponse(a domain.AuthorCommitCount) AuthorCommitCountDto {
	return AuthorCommitCountDto{
		Author:      a.Author,
		CommitCount: a.CommitCount,
	}
}

// AllAuthorCommitCountResponse maps arrary of dto response from array of commit domain entity
func AllAuthorCommitCountResponse(authors []domain.AuthorCommitCount) []AuthorCommitCountDto {
	if len(authors) == 0 {
		return []AuthorCommitCountDto{}
	}

	acResponse := make([]AuthorCommitCountDto, 0, len(authors))

	for _, a := range authors {
		acDto := AuthorCommitCountDto{
			Author:      a.Author,
			CommitCount: a.CommitCount,
		}

		acResponse = append(acResponse, acDto)
	}

	return acResponse
}

// CommitResponse is a mapper of dto commit response from a commit domain entity
func CommitResponse(c domain.Commit) CommitResponseDto {
	return CommitResponseDto{
		CommitID:   c.CommitID,
		Message:    c.Message,
		Author:     c.Author,
		Date:       c.Date,
		URL:        c.URL,
		Repository: c.RepositoryName,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}

// CommitsResponse is a mapper of commit response dto from an array of commit domain entity
func CommitsResponse(commits []domain.Commit) []CommitResponseDto {
	if len(commits) == 0 {
		return []CommitResponseDto{}
	}

	commitsResponse := make([]CommitResponseDto, 0, len(commits))

	for _, c := range commits {
		cr := CommitResponseDto{
			CommitID:   c.CommitID,
			Message:    c.Message,
			Author:     c.Author,
			Date:       c.Date,
			URL:        c.URL,
			Repository: c.RepositoryName,
			CreatedAt:  c.CreatedAt,
			UpdatedAt:  c.UpdatedAt,
		}

		commitsResponse = append(commitsResponse, cr)
	}

	return commitsResponse
}
