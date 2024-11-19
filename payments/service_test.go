package main

import (
	"context"
	"testing"

	"github.com/ashm8206/common/api"
	inmemRegistry "github.com/ashm8206/common/discovery/inmem"
	"github.com/ashm8206/payments/gateway"
	"github.com/ashm8206/payments/processor/inmem"
)

func TestService(t *testing.T) {
	processor := inmem.NewInmem()
	registry := inmemRegistry.NewRegistry()

	gateway := gateway.NewGateway(registry)
	svc := NewService(processor, gateway)
	// svc := NewService(processor)

	t.Run("should create a payment link", func(t *testing.T) {
		link, err := svc.CreatePayment(context.Background(), &api.Order{})
		if err != nil {
			t.Errorf("CreatePayment() error = %v, want nil", err)
		}

		if link == "" {
			t.Error("CreatePayment() link is empty")
		}
	})
}
