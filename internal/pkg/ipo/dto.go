package ipo

type Root struct {
	Props Props `json:"props"`
}

type Props struct {
	PageProps PageProps `json:"pageProps"`
}

type PageProps struct {
	IPOMonthlyList IPOMonthlyList `json:"ipoMonthlyList"`
}

type IPOMonthlyList struct {
	LastMonthCompanies    []Company `json:"lastMonthData"`
	CurrentMonthCompanies []Company `json:"currentMonthData"`
	NextMonthCompanies    []Company `json:"nextMonthData"`
}

type Company struct {
	Name      string `json:"name"`
	Code      string `json:"code"`
	State     string `json:"ipoState"`
	StartDate string `json:"offerSubscriptionStartDate"`
	EndDate   string `json:"offerSubscriptionEndDate"`
	Dart      string `json:"dartLink"`
}
