package presentation

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"payment/app"
	model2 "payment/domain"
	"payment/infra"
)

type HTTP struct {
	payments app.Payments
}

func NewHTTP(p app.Payments) *HTTP {
	return &HTTP{p}
}

func (h *HTTP) Authorize(w http.ResponseWriter, r *http.Request) {
	var req request
	if err := h.decode(r, &req); err != nil {
		h.failed(r, w, err)
		return
	}

	p := h.payment(r)
	m, err := p.Authorize(req.CreditCard, req.Money)
	if err != nil {
		h.failed(r, w, err)
		return
	}

	h.encode(w, response{p.ID(), m})
}

func (h *HTTP) Void(w http.ResponseWriter, r *http.Request) {
	p := h.payment(r)
	m, err := p.Void()
	if err != nil {
		h.failed(r, w, err)
		return
	}

	h.encode(w, response{p.ID(), m})
}

func (h *HTTP) Capture(w http.ResponseWriter, r *http.Request) {
	var req model2.Money
	if err := h.decode(r, &req); err != nil {
		h.failed(r, w, err)
		return
	}

	p := h.payment(r)
	m, err := p.Capture(req)
	if err != nil {
		h.failed(r, w, err)
		return
	}

	h.encode(w, response{p.ID(), m})
}

func (h *HTTP) Refund(w http.ResponseWriter, r *http.Request) {
	var req model2.Money
	if err := h.decode(r, &req); err != nil {
		h.failed(r, w, err)
		return
	}

	p := h.payment(r)
	m, err := p.Refund(req)
	if err != nil {
		h.failed(r, w, err)
		return
	}

	h.encode(w, response{p.ID(), m})
}

func (h *HTTP) payment(r *http.Request) *app.Payment {
	return h.payments.Read(h.id(r), newMerchant(r))
}

func (h *HTTP) id(r *http.Request) model2.ID {
	return model2.NewID(mux.Vars(r)["id"])
}

func (h *HTTP) decode(r *http.Request, d document) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(d); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (h *HTTP) encode(w http.ResponseWriter, d document) {
	json.NewEncoder(w).Encode(d)
}

func (h *HTTP) failed(r *http.Request, w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
	log("ERR %s:%s failed due %s", r.Method, r.URL.String(), err)
}

type request struct {
	CreditCard model2.CreditCard
	Money      model2.Money
}

type response struct {
	ID        model2.ID
	Available model2.Money
}

type document = interface{}

type merchant struct{}

func newMerchant(r *http.Request) *merchant {
	return &merchant{}
}

func (m *merchant) IsAuthenticated() bool {
	return true
}

var log = infra.DefaultLogger.Tag("HTTP").Print
