package vo

type TurnoverReportVO struct {
	DateList     string `json:"dateList"`
	TurnoverList string `json:"turnoverList"`
}

type UserReportVO struct {
	DateList      string `json:"dateList"`
	TotalUserList string `json:"totalUserList"`
	NewUserList   string `json:"newUserList"`
}

type OrderReportVO struct {
	DateList            string  `json:"dateList"`
	OrderCountList      string  `json:"orderCountList"`
	ValidOrderCountList string  `json:"validOrderCountList"`
	TotalOrderCount     int64   `json:"totalOrderCount"`
	ValidOrderCount     int64   `json:"validOrderCount"`
	OrderCompletionRate float64 `json:"orderCompletionRate"`
}

type SalesTop10ReportVO struct {
	NameList   string `json:"nameList"`
	NumberList string `json:"numberList"`
}
