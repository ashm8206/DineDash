package processor

import pb "github.com/ashm8206/common/api"

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}
