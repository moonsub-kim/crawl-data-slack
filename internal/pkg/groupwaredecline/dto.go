package groupwaredecline

type DTO struct {
	UID         string `json:"uid"`
	DocName     string `json:"doc_name"`
	RequestDate string `json:"request_date"` // YYYY-MM-DD
	Drafter     string `json:"drafter"`      // 김멍멍/Dog Kim
	Status      string `json:"status"`       // 반려(김멍멍/Dog Kim)
}
