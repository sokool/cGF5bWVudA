package domain

import (
	"testing"
)

func TestNewMoney(t *testing.T) {
	type (
		have struct {
			amount, currency string
		}

		want error

		case_ struct {
			description string
			have
			want
		}
	)

	scenario := []case_{
		{"no values gives error", have{}, errInvalidMoney},
		{" 0.02 without symbol gives error", have{"2.00", ""}, errInvalidCurrency},
		{"abcd EUR gives error gives error", have{"abcd", "EUR"}, errInvalidMoney},
		{"14.00 alskdjgas gives error", have{"14", "alskdjgas"}, errInvalidCurrency},
		{" 3.99 GBP gives ok", have{"3.99", "gbp"}, nil},
		{" 0.00 USD gives ok", have{" 0.00", "USD"}, nil},
		{"-4.29 PLN gives ok", have{"-4.29", "PLN"}, nil},
		{"88.35 CN¥ gives ok", have{"88.35", "CN¥"}, nil},
	}

	for _, c := range scenario {
		t.Run(c.description, func(t *testing.T) {
			if _, err := NewMoney(c.have.amount, c.have.currency); err != c.want {
				t.Fatalf("expected:%v got:%v", c.want, err)
			}
		})
	}
}
