package groupwaredecline

type DTO struct {
	UID         string `json:"uid"`
	DocName     string `json:"doc_name"`
	RequestDate string `json:"request_date"` // YYYY-MM-DD
	Drafter     string `json:"drafter"`      // 김문섭/Raf Kim
	Status      string `json:"status"`       // 종결(이관우/John Lee)
}
