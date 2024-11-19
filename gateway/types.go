package main

import pb "github.com/ashm8206/common/api"

type CreateOrderRequest struct {
	Order         *pb.Order `"json": order`
	RedirectToURL string    `"json": redirectToURL`
}
