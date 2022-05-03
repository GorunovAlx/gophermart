package withdraw

type WithdrawRepository interface {
	Add(w Withdraw) (int, error)
	GetWithdrawals(userID int) []Withdraw
}
