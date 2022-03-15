package app

import . "payment/domain"

// Payment is a part of application layer.
//
// Has no business logic.
// Checks Merchant identity.
// Read and writes domain objects from databases.
// Transform and returns data to persistent (like HTTP, GRPC, AMQP....) layers.
type Payment struct {
	id           ID
	merchant     Merchant
	transactions Transactions
}

func NewPayment(id ID, m Merchant, t Transactions) *Payment {
	return &Payment{
		id:           id,
		merchant:     m,
		transactions: t,
	}
}

func (t *Payment) ID() ID {
	return t.id
}

func (t *Payment) Authorize(c CreditCard, m Money) (Money, error) {
	if !t.merchant.IsAuthenticated() {
		return Money{}, ErrForbidden
	}

	return t.execute(func(a *Transaction) error { return a.Authorize(c, m) })
}

func (t *Payment) Void() (Money, error) {
	if !t.merchant.IsAuthenticated() {
		return Money{}, ErrForbidden
	}

	return t.execute(func(a *Transaction) error { return a.Void() })
}

func (t *Payment) Capture(m Money) (Money, error) {
	if !t.merchant.IsAuthenticated() {
		return Money{}, ErrForbidden
	}

	return t.execute(func(a *Transaction) error { return a.Capture(m) })
}

func (t *Payment) Refund(m Money) (Money, error) {
	if !t.merchant.IsAuthenticated() {
		return Money{}, ErrForbidden
	}

	return t.execute(func(a *Transaction) error { return a.Refund(m) })
}

func (t *Payment) execute(c command) (available Money, err error) {
	a, err := t.transactions.Read(t.id)
	if err != nil {
		return
	}

	if err = c(a); err != nil {
		return
	}

	if err = t.transactions.Write(a); err != nil {
		return
	}

	return a.Balance(), nil
}

type Payments interface {
	Read(ID, Merchant) *Payment
}

type Transactions interface {
	Read(ID) (*Transaction, error)
	Write(*Transaction) error
}

type command func(*Transaction) error

var ErrForbidden = Err("access forbidden")

type Response struct {
	Transaction ID
	Available   Money
	Error       error
}
