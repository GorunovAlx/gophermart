package withdraw

type WithdrawRepository interface {
	Add(Withdraw)
	GetWithdrawals(userID int) []Withdraw
}
