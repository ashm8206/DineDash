package gateway

import (
	"context"

	pb "github.com/ashm8206/common/api"
)

type StockGateway interface {
	CheckIfItemIsInStock(ctx context.Context, customerID string, items []*pb.ItemsWithQuantity) (bool, []*pb.Item, error)
}
