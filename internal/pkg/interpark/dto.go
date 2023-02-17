package interpark

type Response struct {
	Common struct {
		Message                string `json:"message"`
		RequestURI             string `json:"requestUri"`
		Gtid                   string `json:"gtid"`
		Timestamp              string `json:"timestamp"`
		InternalHTTPStatusCode int    `json:"internalHttpStatusCode"`
	} `json:"common"`
	Data struct {
		RemainSeat []struct {
			PlaySeq       string `json:"playSeq"`
			SeatGrade     string `json:"seatGrade"`
			SeatGradeName string `json:"seatGradeName"`
			RemainCnt     int    `json:"remainCnt"`
		} `json:"remainSeat"`
		Casting interface{} `json:"casting"`
	} `json:"data"`
}
