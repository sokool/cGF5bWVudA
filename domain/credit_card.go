package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CreditCard struct {
	owner
	number
	expiry
	cvv
}

func NewCreditCard(owner, number, expiry, cvv string) (CreditCard, error) {
	var c CreditCard
	var err error

	if c.owner, err = newOwner(owner); err != nil {
		return CreditCard{}, err
	}

	if c.number, err = newNumber(number); err != nil {
		return CreditCard{}, err
	}

	if c.expiry, err = newExpiry(expiry); err != nil {
		return CreditCard{}, err
	}

	if c.cvv, err = newCVV(cvv); err != nil {
		return CreditCard{}, err
	}

	return c, nil
}

func (c CreditCard) IsExpired(d ...time.Time) bool {
	t := time.Now()
	if len(d) != 0 {
		t = d[0]
	}

	return c.expiry.date.Before(t)
}

func (c CreditCard) Number() int {
	return int(c.number)
}

func (c CreditCard) IsZero() bool {
	return c.number == 0
}

func (c CreditCard) String() string {
	return fmt.Sprintf(`%s %d %s %d`,
		c.owner,
		c.number,
		c.expiry,
		c.cvv,
	)
}

func (c CreditCard) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonCreditCard{
		Owner:  string(c.owner),
		Number: c.number.String(),
		Expire: c.expiry.String(),
		CVV:    c.cvv.String(),
	})
}

func (c *CreditCard) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, null) {
		return nil
	}

	var j jsonCreditCard
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	n, err := NewCreditCard(j.Owner, j.Number, j.Expire, j.CVV)
	if err != nil {
		return err
	}

	*c = n
	return nil
}

type owner string

func newOwner(name string) (owner, error) {
	if len(name) < 3 {
		return "", errCreditCardOwner
	}

	return owner(name), nil
}

type number int

func newNumber(num string) (number, error) {
	n, err := strconv.Atoi(strings.ReplaceAll(num, " ", ""))
	if err != nil {
		return 0, errCreditCardNumber
	}

	var m = n / 10
	var s = 0
	for i := 0; m > 0; i++ {
		cur := m % 10
		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		s += cur
		m = m / 10
	}

	if (n%10+s%10)%10 != 0 {
		return 0, errCreditCardNumber
	}

	return number(n), nil
}

func newNumberMust(num string) number {
	n, err := newNumber(num)
	if err != nil {
		panic(err)
	}
	return n
}

func (n number) String() string {
	return strconv.Itoa(int(n))
}

var (
	cardFailures = struct{ auth, capture, refund number }{
		auth:    newNumberMust("4000 0000 0000 0119"),
		capture: newNumberMust("4000 0000 0000 0259"),
		refund:  newNumberMust("4000 0000 0000 3238"),
	}
)

type expiry struct {
	date time.Time
}

func newExpiry(my string) (expiry, error) {
	var e expiry

	if len(my) != 7 {
		return e, errCreditCardExpire
	}

	m, err := strconv.Atoi(my[:2])
	if err != nil {
		return e, errCreditCardExpire
	}

	if m <= 0 || m > 12 {
		return e, errCreditCardExpire
	}

	y, err := strconv.Atoi(my[3:7])
	if err != nil {
		return e, errCreditCardExpire
	}

	return expiry{time.Date(y, time.Month(m)+1, 0, 0, 0, 0, 0, time.UTC)}, nil
}

func (e expiry) String() string {
	return fmt.Sprintf("%02d/%d", e.date.Month(), e.date.Year())
}

type cvv int

func newCVV(code string) (cvv, error) {
	if len(code) != 3 {
		return 0, errCreditCardSecurityCode
	}

	c, err := strconv.Atoi(code)
	if err != nil {
		return 0, errCreditCardSecurityCode
	}

	return cvv(c), nil
}

func (c cvv) String() string {
	return strconv.Itoa(int(c))
}

type jsonCreditCard struct {
	Owner, Number, Expire, CVV string
}

var (
	errCreditCard             = Err("credit card: invalid card data")
	errCreditCardSecurityCode = Err("credit card: invalid cvv code, expected at least 3 digits")
	errCreditCardOwner        = Err("credit card: invalid owner name, expected at least 3 characters")
	errCreditCardNumber       = Err("credit card: invalid number")
	errCreditCardExpire       = Err("credit card: invalid expire date, expected `mm/yyyy` format ")
	errCreditCardExpired      = Err("credit card: expired")
	errCreditCardAuth         = Err("credit card: authorization failure")
	errCreditCardCapture      = Err("credit card: capture failure")
	errCreditCardRefund       = Err("credit card: refund failure")
)
