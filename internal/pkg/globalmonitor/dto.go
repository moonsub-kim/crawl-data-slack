package globalmonitor

import "time"

type DTO struct {
	Nav        map[string]interface{} `json:"nav"`
	ReportList []Report               `json:"reportlist"`
}

type Report struct {
	ID      string `json:"secureId"`
	Company string `json:"auth"`
	Text    string `json:"summary"`
	Title   string `json:"title"`
	Date    string `json:"writeDate"`

	Map map[string]interface{} `json:"-"`
}

//	{
//			"targetPeriodCheck":false
//			"targetPeriodDate":null,
//			"page":1,
//			"pagePerItem":10
//			"sortType":"latest",
//			"lscCd":700,
//			"authLscCd":0,
//			"sscCd":"524910,521090,523050,523070",
//			"authSscCd":0,
//			"cmd":"rl_011",
//			"startDate":"2022-01-16",
//			"endDate":"2022-04-16",
//			"searchItem":"title",
//			"searchStr":"",
//		}
type RequestBody struct {
	TargetPeriodCheck bool    `json:"targetPeriodCheck"`
	TargetPeriodDate  *string `json:"targetPeriodDate,omitempty"`

	Page        int    `json:"page"`
	PagePerItem int    `json:"pagePerItem"`
	SortType    string `json:"sortType"`

	LscCd     int    `json:"lscCd"`
	AuthLscCd int    `json:"authLscCd"`
	SscCd     string `json:"sscCd"`
	AuthSscCd int    `json:"authSscCd"`

	Cmd       string `json:"cmd"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`

	SearchItem string `json:"searchItem"`
	SearchStr  string `json:"searchStr"`
}

func NewRequestBody(startDate time.Time, endDate time.Time) RequestBody {
	return RequestBody{
		TargetPeriodCheck: false,

		Page:        1,
		PagePerItem: 10,
		SortType:    "latest",

		LscCd:     700,
		AuthLscCd: 0,
		SscCd:     "524910,521090,523050,523070",
		AuthSscCd: 0,

		Cmd:       "rl_011",
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),

		SearchItem: "title",
	}
}
