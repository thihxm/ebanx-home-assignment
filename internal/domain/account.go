package domain

type Account struct {
	ID      string `json:"id"`
	Balance int    `json:"balance"`
}

type AccountService interface {
	GetBalance(id string) (int, error)
	Deposit(id string, amount int) (*Account, error)
	Withdraw(id string, amount int) (*Account, error)
	Transfer(originID, destinationID string, amount int) (origin, destination *Account, err error)
	Reset() error
}

type AccountRepository interface {
	FindByID(id string) (*Account, error)
	Upsert(account *Account) (*Account, error)
	Reset() error
}
