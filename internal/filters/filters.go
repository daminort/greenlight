package filters

import (
	"net/url"
	"strings"

	"greenlight.damian.net/internal/queries"
	"greenlight.damian.net/internal/validator"
)

type Filters struct {
	Search   string
	Page     int
	PageSize int
	Sort     string
	Columns  []string
}

type InitParams struct {
	SearchKey   string
	Columns     []string
	SortDefault string
}

type Meta struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	FirstPage    int `json:"first_page"`
	LastPage     int `json:"last_page"`
	TotalRecords int `json:"total_records"`
}

// Filters

func New(values url.Values, params InitParams) *Filters {
	query := queries.New(values)

	f := Filters{
		Search:   query.ReadString(params.SearchKey, ""),
		Page:     query.ReadInt("page", 1),
		PageSize: query.ReadInt("page_size", 10),
		Sort:     query.ReadString("sort", ""),
		Columns:  params.Columns,
	}

	if f.Sort == "" && params.SortDefault != "" {
		f.Sort = params.SortDefault
	}

	return &f
}

func (f *Filters) Validate() validator.ValidationErrors {
	v := validator.New()

	v.Check(validator.GreaterThan(f.Page, 0), "page", "must be greater than zero")
	v.Check(validator.LessThan(f.Page, 10_000_001), "page", "must be a maximum of 10 million")

	v.Check(validator.GreaterThan(f.PageSize, 0), "page_size", "must be greater than zero")
	v.Check(validator.LessThan(f.PageSize, 101), "page_size", "must be a maximum of 100")

	if f.Sort != "" {
		v.Check(validator.InList(f.Sort, f.Columns...), "sort", "invalid sort value")
	}

	return v.Errors
}

func (f *Filters) SortColumn() string {
	for _, column := range f.Columns {
		if f.Sort == column {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	return "title"
}

func (f *Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

func (f *Filters) Limit() int {
	return f.PageSize
}

func (f *Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}

// Meta

func NewMeta(totalRecords, page, pageSize int) *Meta {
	return &Meta{
		CurrentPage:  page,
		PageSize:     pageSize,
		TotalRecords: totalRecords,
		FirstPage:    1,
		LastPage:     (totalRecords + pageSize - 1) / pageSize,
	}
}
