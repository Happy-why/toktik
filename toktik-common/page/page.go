package page

import (
	"github.com/Happy-Why/toktik-common/utils"
	"net/http"
)

// 分页处理

type Page struct {
	DefaultPageSize int64
	MaxPageSize     int64
	PageKey         string // url中page关键字
	PageSizeKey     string // pagesize关键字
}

// InitPage 初始化默认页数大小和最大页数限制以及查询的关键字
func InitPage(defaultPageSize, maxPageSize int64, pageKey, pageSizeKey string) *Page {
	return &Page{
		DefaultPageSize: defaultPageSize,
		MaxPageSize:     maxPageSize,
		PageKey:         pageKey,
		PageSizeKey:     pageSizeKey,
	}
}

// GetPageSizeAndOffset 从请求中获取偏移值和页尺寸
func (p *Page) GetPageSizeAndOffset(r *http.Request) (limit, offset int64) {
	page := utils.StrTo(r.FormValue(p.PageKey)).MustInt64()
	if page <= 0 {
		page = 1
	}
	limit = utils.StrTo(r.FormValue(p.PageSizeKey)).MustInt64()
	if limit <= 0 {
		limit = p.DefaultPageSize
	}
	if limit > p.MaxPageSize {
		limit = p.MaxPageSize
	}
	offset = (page - 1) * limit
	return
}

// CulOffset 计算偏移值
func CulOffset(page, pageSize int32) (offset int32) {
	return (page - 1) * pageSize
}
