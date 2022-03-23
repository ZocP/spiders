package debt

import (
	"qa_spider/pkg/internalsvc/debt/abstract"
)

type Debtor interface {
	GetDebtList(year string, month string) abstract.DebtList
	AddDebtList(debt abstract.Debt) error
	ChangeDebt(ID uint64, status abstract.Status)
}
