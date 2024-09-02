package dtos

import (
	"time"
)

type AllCommitsResponse struct {
	Commits  []CommitResponseDto `json:"commits"`
	PageInfo PagingInfo          `json:"page_info"`
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
