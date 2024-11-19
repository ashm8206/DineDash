package main

import (
	"context"

	pb "github.com/ashm8206/common/api"
)

type PaymentsService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}
