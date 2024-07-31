package request

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func (p *Pagination) Validate() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = 1
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
}
