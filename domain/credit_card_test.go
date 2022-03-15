package domain

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewCreditCard(t *testing.T) {
	c, err := NewCreditCard("Tom", "4000000000000044", "04/2022", "884")

	fmt.Println(c, err)
	fmt.Println(c.IsExpired())

	b, _ := json.MarshalIndent(c, "", "\t")
	fmt.Println(string(b))

	type (
		have struct {
		}

		want struct {
		}

		case_ struct {
			description string
			have
			want
		}
	)

	scenario := []case_{
		{"my first subtest", have{}, want{}},
	}

	check := func(in have) (out want) {
		return
	}

	for _, c := range scenario {
		t.Run(c.description, func(t *testing.T) {
			if out := check(c.have); out != c.want {
				t.Fatalf("expected:%v got:%v", c.want, out)
			}
		})
	}
}
