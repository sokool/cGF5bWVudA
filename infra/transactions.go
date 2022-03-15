package infra

import (
	"pco/app"
	. "pco/domain"
)

type transactions struct{ Events }

func NewTransactions() app.Transactions {
	return &transactions{make(Events)}
}

func (r transactions) Read(id ID) (*Transaction, error) {
	t, err := NewTransaction(id)
	if err != nil {
		return nil, err
	}

	return t, r.read(t)
}

func (r transactions) Write(t *Transaction) error {
	return r.write(t)
}
