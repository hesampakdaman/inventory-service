package domain

type Account struct {
	ID      string
	Owner   string
	Balance float64

	pending []Transaction
}

func NewAccount(ID string, owner string, initialBalance float64) (Account, error) {
	if ID == "" {
		return Account{}, ErrInvalidAccountID
	}
	if owner == "" {
		return Account{}, ErrInvalidOwner
	}
	if initialBalance < 0 {
		return Account{}, ErrNegativeBalance
	}

	return Account{
		ID:      ID,
		Owner:   owner,
		Balance: initialBalance,
	}, nil
}

func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	if err := a.recordTransaction(Deposit, amount); err != nil {
		return err
	}

	a.Balance += amount

	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	if amount > a.Balance {
		return ErrInsufficientFunds
	}

	if err := a.recordTransaction(Withdrawal, amount); err != nil {
		return err
	}

	a.Balance -= amount

	return nil
}

func (a *Account) Transfer(to *Account, amount float64) error {
	if a.ID == to.ID {
		return ErrSelfTransfer
	}

	if err := a.Withdraw(amount); err != nil {
		return err
	}

	if err := to.Deposit(amount); err != nil {
		return err
	}

	return nil
}

func (a *Account) recordTransaction(ttype TransactionType, amount float64) error {
	txn, err := NewTransaction(a.ID, ttype, amount)
	if err != nil {
		return err
	}

	a.pending = append(a.pending, txn)

	return err
}
