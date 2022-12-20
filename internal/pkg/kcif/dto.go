package kcif

type DTO struct {
	ID      string
	Date    string
	Title   string
	pdfURL  string
	Content string
}

type PageRequest struct {
	IntReportID   string `url:"intReportID"`
	CurrentPage   string `url:"currentPage"`
	RowsPerPage   string `url:"rowsPerPage"`
	IntSection1   string `url:"intSection1"`
	IntSection2   string `url:"intSection2"`
	IntBoardID    string `url:"intBoardID"`
	Regular       string `url:"regular"`
	AnalysisBrief string `url:"AnalysisBrief"`
	Intperiod1    string `url:"intperiod1"`
	Intperiod2    string `url:"intperiod2"`
	OrderValue    string `url:"orderValue"`
	S_title       string `url:"s_title"`
	S_word        string `url:"s_word"`
}

func newPageRequest(id string) PageRequest {
	return PageRequest{
		IntReportID:   id,
		CurrentPage:   "1",
		RowsPerPage:   "15",
		IntSection1:   "1",
		IntSection2:   "5",
		IntBoardID:    "5",
		Regular:       "",
		AnalysisBrief: "",
		Intperiod1:    "",
		Intperiod2:    "",
		OrderValue:    "",
		S_title:       "true",
		S_word:        "",
	}
}
