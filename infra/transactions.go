package infra

import (
	"payment/app"
	"payment/domain"
)

type transactions struct{ Events }

func NewTransactions() app.Transactions {
	return &transactions{make(Events)}
}

func (r transactions) Read(id domain.ID) (*domain.Transaction, error) {
	t, err := domain.NewTransaction(id)
	if err != nil {
		return nil, err
	}

	return t, r.read(t)
}

func (r transactions) Write(t *domain.Transaction) error {
	return r.write(t)
}
