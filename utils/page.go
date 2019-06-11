package utils

type Page struct {
	PageSize    int         `json:"-"`
	TotalCount  int         `json:"total_count"`
	TotalPage   int         `json:"-"`
	CurrentPage int         `json:"-"`
	StartIndex  int         `json:"-"`
	Data        interface{} `json:"data"`
}

func NewPage(currentPage int, pageSize int, totalCount int) *Page {
	page := Page{}
	if pageSize == 0 {
		page.PageSize = 20
	}
	if currentPage == 0 {
		page.CurrentPage = 1
	}
	page.PageSize = pageSize
	page.CurrentPage = currentPage
	page.StartIndex = (page.CurrentPage - 1) * page.PageSize
	page.TotalCount = totalCount
	if page.TotalCount%page.PageSize == 0 {
		page.TotalPage = page.TotalCount / page.PageSize
	} else {
		page.TotalPage = page.TotalCount/page.PageSize + 1
	}
	return &page
}
