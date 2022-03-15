package domain

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid"
)

// Transaction represents business logic of Money flow and management.
//
// It holds business invariants and checks them on each method call in order
// to deliver consistency for customer.
//
// I did assumption that:
// CreditCard can be stored in payment gateway.
// Payment Gateway itself manages bank accounts
type Transaction struct {
	id         ID
	card       CreditCard
	authorized Money
	balance    Money
	voided     bool

	uncommitted []Event
}

func NewTransaction(id ID) (*Transaction, error) {
	return &Transaction{id: id}, nil
}

func (a *Transaction) ID() string {
	return string(a.id)
}

func (a *Transaction) Authorize(c CreditCard, m Money) error {
	switch {
	case !a.authorized.IsZero():
		return errTxAlreadyAuthorized
	case !m.IsPositive():
		return errInsufficientAmount
	case c.IsZero():
		return errCreditCard
	case c.IsExpired():
		return errCreditCardExpired
	case c.number == cardFailures.auth:
		return errCreditCardAuth
	}

	return a.append(TransactionAuthorized{c, m})
}

func (a *Transaction) Void() error {
	switch {
	case a.voided:
		return nil
	case a.authorized.IsZero():
		return errTxNotFound
	case !a.authorized.sub(a.balance).IsZero():
		return errTxVoidRejected
	}

	return a.append(TransactionVoided{})
}

func (a *Transaction) Capture(m Money) error {
	switch {
	case a.authorized.IsZero():
		return errTxNotFound
	case a.card.number == cardFailures.capture:
		return errCreditCardCapture
	case a.voided:
		return errTxVoided
	case a.balance.lower(m):
		return errTxCaptureExceeded
	case !m.IsPositive():
		return errInsufficientAmount
	}

	return a.append(TransactionCaptured{m})
}

func (a *Transaction) Refund(m Money) error {
	switch {
	case a.authorized.IsZero():
		return errTxNotFound
	case a.card.number == cardFailures.refund:
		return errCreditCardRefund
	case a.voided:
		return errTxVoided
	case a.authorized.sub(a.balance).lower(m):
		return errTxRefundExceeded
	case !m.IsPositive():
		return errInsufficientAmount
	}

	return a.append(TransactionRefunded{m})
}

func (a *Transaction) Balance() Money {
	return a.balance
}

func (a *Transaction) Commit(e Event, at time.Time) error {
	switch e := e.(type) {
	case TransactionAuthorized:
		a.authorized, a.balance, a.card = e.Money, e.Money, e.CreditCard
	case TransactionCaptured:
		a.balance = a.balance.sub(e.Money)
	case TransactionRefunded:
		a.balance = a.balance.add(e.Money)
	case TransactionVoided:
		a.voided = true
	}

	return nil
}

func (a *Transaction) Uncommitted(clear bool) []Event {
	defer func() {
		if clear {
			a.uncommitted = []Event{}
		}
	}()

	return a.uncommitted
}

func (a *Transaction) append(events ...Event) error {
	a.uncommitted = append(a.uncommitted, events...)
	return nil
}

var (
	errTxNotFound          = Err("transaction: not found")
	errTxAlreadyAuthorized = Err("transaction: already authorized")
	errTxVoidRejected      = Err("transaction: void rejected")
	errTxCaptureExceeded   = Err("transaction: capture amount exceeded")
	errTxRefundExceeded    = Err("transaction: refund amount exceeded")
	errTxVoided            = Err("transaction: voided")
)

type ID string

func NewID(s ...string) ID {
	if len(s) != 0 && s[0] != "" {
		return ID(s[0])
	}

	return ID(gonanoid.MustID(20))
}

type (
	Event = interface{}

	TransactionAuthorized struct {
		CreditCard
		Money
	}

	TransactionVoided struct {
	}

	TransactionCaptured struct {
		Money
	}

	TransactionRefunded struct {
		Money
	}
)
