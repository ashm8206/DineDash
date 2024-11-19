package main

import (
	"context"

	pb "github.com/ashm8206/common/api"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DbName   = "orders"
	CollName = "orders"
)

type store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *store {
	return &store{db}
}

func (s *store) Create(ctx context.Context, o Order) (primitive.ObjectID, error) {
	col := s.db.Database(DbName).Collection(CollName)

	newOrder, err := col.InsertOne(ctx, o)

	id := newOrder.InsertedID.(primitive.ObjectID)
	return id, err
}

func (s *store) Get(ctx context.Context, id, customerID string) (*Order, error) {
	col := s.db.Database(DbName).Collection(CollName)

	oID, _ := primitive.ObjectIDFromHex(id)

	var o Order
	err := col.FindOne(ctx, bson.M{
		"_id":        oID,
		"customerID": customerID,
	}).Decode(&o)

	return &o, err
}

func (s *store) Update(ctx context.Context, id string, newOrder *pb.Order) error {
	col := s.db.Database(DbName).Collection(CollName)

	oID, _ := primitive.ObjectIDFromHex(id)

	_, err := col.UpdateOne(ctx,
		bson.M{"_id": oID},
		bson.M{"$set": bson.M{
			"paymentLink": newOrder.PaymentLink,
			"status":      newOrder.Status,
		}})

	return err
}

//Temporary w/o persistence

// var orders = make([]*pb.Order, 0)

// type store struct {
// 	// add MongoDB
// }
// func NewStore() *store {
// 	return &store{}
// }

// func (s *store) Create(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Item) (string, error) {
// 	id := "42"
// 	orders = append(orders, &pb.Order{
// 		ID:          id,
// 		CustomerID:  p.CustomerID,
// 		Status:      "pending",
// 		Items:       items,
// 		PaymentLink: "",
// 	})
// 	return id, nil
// }

// func (s *store) Update(ctx context.Context, id string, newOrder *pb.Order) error {

// 	for i, order := range orders {
// 		if order.ID == id {
// 			orders[i].Status = newOrder.Status
// 			orders[i].PaymentLink = newOrder.PaymentLink
// 			return nil
// 		}
// 	}

// 	// col := s.db.Database(DbName).Collection(CollName)

// 	// oID, _ := primitive.ObjectIDFromHex(id)

// 	// var o Order
// 	// err := col.FindOne(ctx, bson.M{
// 	// 	"_id":        oID,
// 	// 	"customerID": customerID,
// 	// }).Decode(&o)

// 	return errors.New("no orders found")
// }

// func (s *store) Get(ctx context.Context, id, customerID string) (*pb.Order, error) {

// 	for _, o := range orders {
// 		if o.ID == id && o.CustomerID == customerID {
// 			return o, nil
// 		}
// 	}

// 	// col := s.db.Database(DbName).Collection(CollName)

// 	// oID, _ := primitive.ObjectIDFromHex(id)

// 	// var o Order
// 	// err := col.FindOne(ctx, bson.M{
// 	// 	"_id":        oID,
// 	// 	"customerID": customerID,
// 	// }).Decode(&o)

// 	return nil, errors.New("no orders found")
// }
