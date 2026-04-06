package pagination

type Page struct {
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
}

type Config struct {
	DefaultPage    int
	DefaultPerPage int
	MinPerPage     int
	MaxPerPage     int
}

var DefaultConfig = Config{
	DefaultPage:    1,
	DefaultPerPage: 20,
	MinPerPage:     1,
	MaxPerPage:     100,
}

// Normalize applies sane defaults when Page or PerPage are out of range.
// Out-of-range values are silently replaced with defaults — callers should
// not rely on specific values being preserved.
func (p *Page) Normalize(cfg Config) {
	if p.Page < 1 {
		p.Page = cfg.DefaultPage
	}
	if p.PerPage < cfg.MinPerPage || p.PerPage > cfg.MaxPerPage {
		p.PerPage = cfg.DefaultPerPage
	}
}

func (p Page) Offset() int {
	return (p.Page - 1) * p.PerPage
}
