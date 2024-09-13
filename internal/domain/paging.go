package domain

import "github.com/kenmobility/git-api-service/internal/http/dtos"

type (
	APIPaging struct {
		Limit     int
		Page      int
		Sort      string
		Direction string
	}

	PagingInfo struct {
		TotalCount  int64
		Page        int
		HasNextPage bool
		Count       int
	}
)

func (p APIPaging) ToDto() dtos.APIPagingDto {
	return dtos.APIPagingDto{
		Limit:     p.Limit,
		Page:      p.Page,
		Sort:      p.Sort,
		Direction: p.Direction,
	}
}

func (p PagingInfo) ToDto() dtos.PagingInfo {
	return dtos.PagingInfo{
		TotalCount:  p.TotalCount,
		Page:        p.Page,
		HasNextPage: p.HasNextPage,
		Count:       p.Count,
	}
}

func FromDtoPaging(query dtos.APIPagingDto) APIPaging {
	return APIPaging{
		Limit:     query.Limit,
		Page:      query.Page,
		Sort:      query.Sort,
		Direction: query.Direction,
	}
}
