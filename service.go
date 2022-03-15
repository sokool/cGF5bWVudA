package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"pco/app"
	"pco/domain"
	"pco/infra"
	"pco/presentation"
)

type Service struct {
	transaction app.Transactions
}

func NewService() *Service {
	return &Service{
		transaction: infra.NewTransactions(),
	}
}

func (s *Service) Read(id domain.ID, m app.Merchant) *app.Payment {
	return app.NewPayment(id, m, s.transaction)
}

func (s *Service) Run() error {
	h := presentation.NewHTTP(s)
	r := mux.NewRouter()
	r.HandleFunc("/transactions/authorize", h.Authorize).Methods("POST")
	r.HandleFunc("/transactions/{id}/void", h.Void).Methods("PUT")
	r.HandleFunc("/transactions/{id}/capture", h.Capture).Methods("PUT")
	r.HandleFunc("/transactions/{id}/refund", h.Refund).Methods("PUT")

	return http.ListenAndServe("", r)
}
