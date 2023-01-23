package navercareer

import (
	"encoding/json"
	"time"
)

type DTO struct {
	Title string
	Info  string
	URL   string
}

type Request struct {
	Where     string `url:"where"`
	Query     string `url:"query"`
	NLUQuery  string `url:"nlu_query"`
	APIType   string `url:"api_type"`
	SortOrder string `url:"so"`
	Start     string `url:"start"`
	Callback  string `url:"_callback"`
	Timestamp int64  `url:"_"`
}

func newPageRequest(query string) (Request, error) {
	nluQuery, err := json.Marshal(map[string]string{
		"q":           query,
		"unknownType": query,
	})
	if err != nil {
		return Request{}, err
	}

	return Request{
		Where:     "pc_bridge_list",
		Query:     query,
		NLUQuery:  string(nluQuery),
		APIType:   "1",
		SortOrder: "elapsed_time.asc",
		Start:     "1",
		Timestamp: time.Now().Unix(),
	}, nil
}
