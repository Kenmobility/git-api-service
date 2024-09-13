package domain

import (
	"time"

	"github.com/kenmobility/git-api-service/internal/http/dtos"
)

type Commit struct {
	CommitID       string
	Message        string
	Author         string
	Date           time.Time
	URL            string
	RepositoryName string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (c Commit) ToDto() dtos.CommitResponseDto {
	return dtos.CommitResponseDto{
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

func CommitsResponse(commits []Commit) []dtos.CommitResponseDto {
	if len(commits) == 0 {
		return []dtos.CommitResponseDto{}
	}

	commitsResponse := make([]dtos.CommitResponseDto, 0, len(commits))

	for _, c := range commits {
		cr := dtos.CommitResponseDto{
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

type AuthorCommitCount struct {
	Author      string
	CommitCount int
}

func (a AuthorCommitCount) ToDto() dtos.AuthorCommitCountDto {
	return dtos.AuthorCommitCountDto{
		Author:      a.Author,
		CommitCount: a.CommitCount,
	}
}

func AuthorsCommitCountResponse(authors []AuthorCommitCount) []dtos.AuthorCommitCountDto {
	if len(authors) == 0 {
		return []dtos.AuthorCommitCountDto{}
	}

	acResponse := make([]dtos.AuthorCommitCountDto, 0, len(authors))

	for _, a := range authors {
		acDto := dtos.AuthorCommitCountDto{
			Author:      a.Author,
			CommitCount: a.CommitCount,
		}

		acResponse = append(acResponse, acDto)
	}

	return acResponse
}
