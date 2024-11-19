package inmem

import pb "github.com/ashm8206/common/api"

type Inmem struct{}

func NewInmem() *Inmem {
	return &Inmem{}
}

func (i *Inmem) CreatePaymentLink(*pb.Order) (string, error) {
	return "dummy-link", nil
}
