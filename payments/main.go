package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/ashm8206/common"
	"github.com/ashm8206/common/broker"
	"github.com/ashm8206/common/discovery"
	"github.com/ashm8206/common/discovery/consul"
	"github.com/ashm8206/payments/gateway"
	stripeProcessor "github.com/ashm8206/payments/processor/stripe"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stripe/stripe-go/v81"
	"google.golang.org/grpc"
)

var (
	serviceName          = "orders"
	grpcAddr             = common.EnvString("GRPC_ADDR", "localhost:2001")
	consulAddr           = common.EnvString("CONSUL_ADDR", "localhost:8500")
	amqpUser             = common.EnvString("RABBITMQ_USER", "guest")
	amqpPass             = common.EnvString("RABBITMQ_PASS", "guest")
	amqpHost             = common.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort             = common.EnvString("RABBITMQ_PORT", "5672")
	stripeKey            = common.EnvString("STRIPE_KEY", "")
	httpAddr             = common.EnvString("HTTP_ADDR", "localhost:8081")
	endpointStripeSecret = common.EnvString("ENDPOINT_STRIPE_SECRET", "whsec..")
	jaegerAddr           = common.EnvString("JAEGER_ADDR", "localhost:4318")
)

func main() {

	if err := common.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr); err != nil {
		log.Fatal("failed to set global tracer")
	}

	// registry Consul
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("failed to health check")
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	stripe.Key = stripeKey

	// Broker Connection
	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	stripeProcessor := stripeProcessor.NewProcessor()
	gateway := gateway.NewGateway(registry)
	svc := NewService(stripeProcessor, gateway)

	ampqConsumer := NewConsumer(svc)

	go ampqConsumer.Listen(ch)

	// http server
	mux := http.NewServeMux()

	httpServer := NewPaymentHTTPHandler(ch)
	httpServer.registerRoutes(mux)

	go func() {
		log.Printf("Starting HTTP server at %s", httpAddr)
		if err := http.ListenAndServe(httpAddr, mux); err != nil {
			log.Fatal("failed to start http server")
		}
	}()

	// grpc Server
	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)

	if err != nil {
		log.Fatalf("failed to listen %v", err)

	}

	defer l.Close()

	log.Println("GRPC Server started @", grpcAddr)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}

}
