package quasarzone

type DTO struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Category  string `json:"category"`
	PriceInfo string `json:"price_info"`
	Date      string `json:"date"`
}

func (d DTO) isEmpty() bool {
	return d == DTO{}
}
