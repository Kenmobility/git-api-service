package dtos

import "github.com/kenmobility/git-api-service/internal/domain"

type (
	APIPagingDto struct {
		Limit     int    `json:"limit,omitempty"`
		Page      int    `json:"page,omitempty"`
		Sort      string `json:"sort,omitempty"`
		Direction string `json:"direction,omitempty"`
	}

	PagingInfoDto struct {
		TotalCount  int64 `json:"totalCount"`
		Page        int   `json:"page"`
		HasNextPage bool  `json:"hasNextPage"`
		Count       int   `json:"count"`
	}
)

// PagingInfoResponse is a mapper to pagingInfo dto from domain entity
func PagingInfoResponse(p domain.PagingInfo) PagingInfoDto {
	return PagingInfoDto{
		TotalCount:  p.TotalCount,
		Page:        p.Page,
		HasNextPage: p.HasNextPage,
		Count:       p.Count,
	}
}

// PagingDataFromPagingDto is a mapper from APIPagingDto to domain entity APIPagingData
func PagingDataFromPagingDto(query APIPagingDto) domain.APIPagingData {
	return domain.APIPagingData{
		Limit:     query.Limit,
		Page:      query.Page,
		Sort:      query.Sort,
		Direction: query.Direction,
	}
}
