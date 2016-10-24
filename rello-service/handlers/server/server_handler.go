package handler

// This file contains the Service definition, and a basic service
// implementation. It also includes service middlewares.

import (
	_ "errors"
	_ "time"

	"golang.org/x/net/context"

	_ "github.com/go-kit/kit/log"
	_ "github.com/go-kit/kit/metrics"

	pb "github.com/adamryman/rello/rello-service"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() Service {
	return relloService{}
}

type relloService struct{}

// CheckListWebhook implements Service.
func (s relloService) CheckListWebhook(ctx context.Context, in *pb.ChecklistUpdate) (*pb.Empty, error) {
	_ = ctx
	_ = in

	response := pb.Empty{}
	return &response, nil
}

type Service interface {
	CheckListWebhook(ctx context.Context, in *pb.ChecklistUpdate) (*pb.Empty, error)
}
