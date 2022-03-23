package debt

type RDebt struct {
	Title       string   `json:"title"`
	Members     []string `json:"members"`
	IsOwe       bool     `json:"is_owe"`
	OweTime     string   `json:"owe_time"`
	RevertsTime string   `json:"reverts_time"`
}

type RDebtList struct {
	Year     string  `json:"year"`
	Month    string  `json:"month"`
	DebtList []RDebt `json:"debt_list"`
}
