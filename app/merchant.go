package app

type Merchant interface {
	IsAuthenticated() bool
}
