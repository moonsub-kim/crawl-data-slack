package financialreport

type DTO struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Text    string `json:"text"`
	Type    string `json:"type"`
	PDFURL  string
}
