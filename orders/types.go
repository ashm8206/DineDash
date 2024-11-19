package main

import (
	"context"

	pb "github.com/ashm8206/common/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrdersService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest, []*pb.Item) (*pb.Order, error)
	ValidateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
	GetOrder(context.Context, *pb.GetOrderRequest) (*pb.Order, error)
	UpdateOrder(context.Context, *pb.Order) (*pb.Order, error)
}

type OrdersStore interface {
	Create(context.Context, Order) (primitive.ObjectID, error)
	Get(ctx context.Context, id, customerID string) (*Order, error)
	//Temp without persistence
	// Create(context.Context, *pb.CreateOrderRequest, []*pb.Item) (string, error)
	// Get(ctx context.Context, id, customerID string) (*pb.Order, error)

	Update(ctx context.Context, id string, o *pb.Order) error
}

// Model Instance of the Order

type Order struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	CustomerID  string             `bson:"customerID,omitempty"`
	Status      string             `bson:"status,omitempty"`
	PaymentLink string             `bson:"paymentLink,omitempty"`
	Items       []*pb.Item         `bson:"items,omitempty"`
}

func (o *Order) ToProto() *pb.Order {
	return &pb.Order{
		ID:          o.ID.Hex(),
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
	}
}
