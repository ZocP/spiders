package abstract

type Debt struct {
	ID          uint64   `json:"id"`
	Title       string   `json:"title"`
	Members     []string `json:"members"`
	IsOwe       bool     `json:"is_owe"`
	OweTime     string   `json:"owe_time"`
	RevertsTime string   `json:"reverts_time"`
}

type DebtList struct {
	Year     string `json:"year"`
	Month    string `json:"month"`
	DebtList []Debt `json:"debt_list"`
}

type Status struct {
	Title       string
	Members     []string
	IsOwe       bool
	OweTime     string
	RevertsTime string
}
