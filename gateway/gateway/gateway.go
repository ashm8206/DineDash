package gateway

import (
	"context"

	pb "github.com/ashm8206/common/api"
)

type OrdersGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	GetOrder(ctx context.Context, orderID, customerID string) (*pb.Order, error)
}
