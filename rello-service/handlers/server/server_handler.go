package handler

// This file contains the Service definition, and a basic service
// implementation. It also includes service middlewares.

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/adamryman/rello"
	pb "github.com/adamryman/rello/rello-service"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.RelloServer {
	return relloService{}
}

type relloService struct{}

// CheckListWebhook implements Service.
func (s relloService) CheckListWebhook(ctx context.Context, in *pb.ChecklistUpdate) (*pb.Empty, error) {
	hash := ctx.Value("X-Trello-Webhook")
	// TODO: verify this hash https://developers.trello.com/apis/webhooks
	_ = hash
	fmt.Println("X-Trello-Webhook: ", hash)

	rello.HandleUpdate(in)

	return &pb.Empty{}, nil
}

type Service interface {
	CheckListWebhook(ctx context.Context, in *pb.ChecklistUpdate) (*pb.Empty, error)
}
