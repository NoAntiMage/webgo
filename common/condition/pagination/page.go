package pagination

type Page struct {
	PageNum  int `form:"pageNum" default:"1" validate:"gt=0"`
	PageSize int `form:"pageSize" default:"10" validate:"gt=0"`
}

func (p *Page) ToSqlLimit() (limit int, offset int) {
	return p.PageSize, (p.PageNum - 1) * p.PageSize

}
