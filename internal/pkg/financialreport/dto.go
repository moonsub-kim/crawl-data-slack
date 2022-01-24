package financialreport

type DTO struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Text    string `json:"text"`
	Company string `json:"company"`
	PDFURL  string
}
