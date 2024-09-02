package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kenmobility/git-api-service/internal/http/dtos"
)

func getPagingInfo(c *gin.Context) dtos.APIPagingDto {
	var paging dtos.APIPagingDto

	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	sort := c.Query("sort")
	direction := c.Query("direction")

	paging.Limit = limit
	paging.Page = page
	paging.Sort = sort
	paging.Direction = direction

	return paging
}
