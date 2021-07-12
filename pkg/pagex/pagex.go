package pagex

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var defaultPage = 1
var defaultPageSize = 10

// PaginateResult paginate result
type PaginateResult struct {
	Page       int                      `json:"page"`
	PageSize   int                      `json:"page_size"`
	TotalPage  int                      `json:"total_page"`
	TotalCount int64                    `json:"total_count"`
	List       []map[string]interface{} `json:"list"`
}

// Paginator paginator
type Paginator struct {
	Page     int
	PageSize int
}

// DefaultPaginator new paginator with default config (page:1, pageSize:10)
func DefaultPaginator() *Paginator {
	return &Paginator{Page: defaultPage, PageSize: defaultPageSize}
}

// NewPaginator new paginator with custom config
func NewPaginator(page, pageSize int) *Paginator {
	return &Paginator{Page: page, PageSize: pageSize}
}

// GinPaginator new paginator with gin context
func GinPaginator(ctx *gin.Context) *Paginator {
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	return &Paginator{Page: page, PageSize: pageSize}
}

// Paginate do paginate
func (p *Paginator) Paginate(query *gorm.DB) PaginateResult {
	if p.Page <= 0 {
		p.Page = defaultPage
	}
	if p.PageSize <= 0 {
		p.PageSize = defaultPageSize
	}
	result := PaginateResult{}
	result.List = []map[string]interface{}{}

	total := p.TotalCount(query)
	offset := (p.Page - 1) * p.PageSize
	query.Offset(offset).Limit(p.PageSize).Find(&result.List)

	result.Page = p.Page
	result.PageSize = p.PageSize
	result.TotalPage = int(math.Ceil(float64(total) / float64(p.PageSize)))
	result.TotalCount = total

	return result
}

// TotalCount get total data count
func (p *Paginator) TotalCount(query *gorm.DB) int64 {
	var count int64
	query.Count(&count)
	return count
}
