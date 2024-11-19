package main

import (
	"context"
	"log"

	common "github.com/ashm8206/common"
	pb "github.com/ashm8206/common/api"
	"github.com/ashm8206/orders/gateway"
)

type service struct {
	store   OrdersStore
	gateway gateway.StockGateway
}

func NewService(store OrdersStore, gateway gateway.StockGateway) *service {
	return &service{store, gateway}
}

// func (UnimplementedOrderServiceServer) GetOrder(context.Context, *GetOrderRequest) (*Order, error) {
// 	return nil, status.Errorf(codes.Unimplemented, "method GetOrder not implemented")
// }

func (s *service) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {

	o, err := s.store.Get(ctx, p.OrderID, p.CustomerID)
	if err != nil {
		return nil, err
	}

	return o.ToProto(), nil

}

func (s *service) UpdateOrder(ctx context.Context, o *pb.Order) (*pb.Order, error) {
	err := s.store.Update(ctx, o.ID, o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (s *service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {

	id, err := s.store.Create(ctx, Order{
		CustomerID:  p.CustomerID,
		Status:      "pending",
		Items:       items,
		PaymentLink: "",
	})

	// Without persistence
	// id, err := s.store.Create(ctx, p, items)

	if err != nil {
		return nil, err
	}

	o := &pb.Order{
		ID: id.Hex(),
		// ID:         id,
		CustomerID: p.CustomerID,
		Status:     "Pending",
		Items:      items,
	}

	return o, nil
}

func (s *service) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item, error) {
	if len(p.Items) == 0 {
		return nil, common.ErrNoItems
	}

	mergedItems := mergeItemsQuantities(p.Items)
	log.Printf("Items were Merged and Validated Succesfully")
	log.Println(mergedItems)

	//Temporary

	// var itemsWithPrice []*pb.Item
	// for _, i := range mergedItems {
	// 	itemsWithPrice = append(itemsWithPrice, &pb.Item{
	// 		PriceID:  "price_1QJjMOP4O7X0edOL03efPdaM",
	// 		ID:       i.ID,
	// 		Quantity: i.Quantity,
	// 	})
	// }
	// validate with the stock service

	inStock, items, err := s.gateway.CheckIfItemIsInStock(ctx, p.CustomerID, mergedItems)
	if err != nil {
		return nil, err
	}
	if !inStock {
		return items, common.ErrNoStock
	}
	log.Println(items)
	return items, nil
}

func mergeItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	merged := make([]*pb.ItemsWithQuantity, 0)

	for _, item := range items {
		found := false
		for _, finalItem := range merged {
			if finalItem.ID == item.ID {
				finalItem.Quantity += item.Quantity
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, item)
		}
	}

	return merged
}
