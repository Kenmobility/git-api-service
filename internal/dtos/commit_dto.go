package dtos

import (
	"github.com/kenmobility/github-api-service/internal/domains/models"
)

type AllCommitsResponse struct {
	Commits  []models.Commit `json:"commits"`
	PageInfo PagingInfo      `json:"page_info"`
}
