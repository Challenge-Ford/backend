package pagination

type Meta struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

func NewMeta(p Page, total int) Meta {
	totalPages := (total + p.PerPage - 1) / p.PerPage
	return Meta{
		Page:       p.Page,
		PerPage:    p.PerPage,
		Total:      total,
		TotalPages: totalPages,
	}
}

type Result[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
}

func NewResult[T any](data []T, p Page, total int) Result[T] {
	return Result[T]{
		Data: data,
		Meta: NewMeta(p, total),
	}
}
