package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pb "github.com/ashm8206/common/api"
	"github.com/ashm8206/common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrdersService
	channel *amqp.Channel
}

func NewGRPCHandler(grpcServer *grpc.Server, service OrdersService, channel *amqp.Channel) {
	handler := &grpcHandler{
		service: service,
		channel: channel,
	}
	pb.RegisterOrderServiceServer(grpcServer, handler)

}

func (h *grpcHandler) UpdateOrder(ctx context.Context, p *pb.Order) (*pb.Order, error) {
	return h.service.UpdateOrder(ctx, p)

}

func (h *grpcHandler) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	return h.service.GetOrder(ctx, p)

}

func (h *grpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New Order Recieved Order %v", p)


	q, err := h.channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)

	}

	// items, err := h.service.ValidateOrder(ctx, p)
	// if err != nil {
	// 	return nil, err
	// }

	// o, err := h.service.CreateOrder(ctx, p, items)

	// if err != nil {
	// 	return nil, err
	// }

	tr := otel.Tracer("amqp")
	amqpContext, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - publish - %s", q.Name))
	defer messageSpan.End()

	items, err := h.service.ValidateOrder(amqpContext, p)
	if err != nil {
		return nil, err
	}

	o, err := h.service.CreateOrder(amqpContext, p, items)
	if err != nil {
		return nil, err
	}


	marshalledOrder, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	// inject the headers
	headers := broker.InjectAMQPHeaders(amqpContext)

	h.channel.PublishWithContext(amqpContext, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         marshalledOrder,
		DeliveryMode: amqp.Persistent,
		Headers:      headers,
	})

	// //temp
	// h.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
	// 	ContentType:  "application/json",
	// 	Body:         marshalledOrder,
	// 	DeliveryMode: amqp.Persistent,
	// 	// Headers:      headers,
	// })

	return o, nil
}
