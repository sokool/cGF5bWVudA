package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Money
//
// todo check currency when operations on multiple Money is performed
type Money struct {
	amount // todo find better solution as value object instead pointer
	currency
}

func NewMoney(s, currency string) (Money, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return Money{}, errInvalidMoney
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return Money{}, errInvalidMoney
	}

	c, err := newCurrency(currency)
	if err != nil {
		return Money{}, err
	}

	return Money{amount(f), c}, nil
}

func (m Money) greater(than Money) bool {
	return m.Amount() > than.Amount()
}

func (m Money) lower(than Money) bool {
	return m.Amount() < than.Amount()
}

func (m Money) add(n Money) Money {
	a := amount(m.Amount() + n.Amount())
	return Money{a, m.currency}
}

func (m Money) sub(n Money) Money {
	a := amount(m.Amount() - n.Amount())
	return Money{a, m.currency}
}

func (m Money) IsPositive() bool {
	if m.IsZero() {
		return false
	}
	return m.Amount() > 0
}

func (m Money) IsZero() bool {
	return m.Amount() == 0
}

func (m Money) String() string {
	return fmt.Sprintf("%s%s", m.amount, m.currency)
}

func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonMoney{
		Amount:   m.amount.String(),
		Currency: m.currency.Symbol(),
	})
}

func (m *Money) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, null) {
		return nil
	}

	var j jsonMoney
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	v, err := NewMoney(j.Amount, j.Currency)
	if err != nil {
		return err
	}

	*m = v
	return nil
}

type amount float64

func (p amount) Pennies() int {
	return int(p * 100)
}

func (p amount) Amount() float64 {
	return float64(p)
}

func (p amount) String() string {
	return fmt.Sprintf("%.2f", p.Amount())
}

type currency string

func newCurrency(symbol string) (currency, error) {
	if n := len([]rune(symbol)); n == 0 || n > 3 {
		return "", errInvalidCurrency
	}

	return currency(strings.ToTitle(symbol)), nil
}

func (p currency) Symbol() string {
	return string(p)
}

func (p currency) String() string {
	return string(p)
}

type jsonMoney struct {
	Amount, Currency string
}

var (
	errInvalidMoney       = Err("money: invalid amount, expected ie 149.99 format")
	errInvalidCurrency    = Err("money: invalid symbol, ie USD is expected")
	errInsufficientAmount = Err("money: insufficient amount")
)
