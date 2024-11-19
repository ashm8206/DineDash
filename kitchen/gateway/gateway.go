package gateway

import (
	"context"

	pb "github.com/ashm8206/common/api"
)

type KitchenGateway interface {
	UpdateOrder(context.Context, *pb.Order) error
}
