package goo

type Pagination struct {
	Num     int    `json:"page_num"`
	Size    int    `json:"page_size"`
	OrderBy string `json:"order_by"`
}

func (p *Pagination) Limit() int {
	if p.Size == 0 {
		return 20
	}
	return p.Size
}

func (p *Pagination) Offset() int {
	return (p.Num - 1) * p.Limit()
}
