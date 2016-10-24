package svc

// This file provides server-side bindings for the gRPC transport.
// It utilizes the transport/grpc.Server.

import (
	//stdopentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"

	//"github.com/go-kit/kit/log"
	//"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	pb "github.com/adamryman/rello/rello-service"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC AddServer.
func MakeGRPCServer(ctx context.Context, endpoints Endpoints /*, tracer stdopentracing.Tracer, logger log.Logger*/) pb.RelloServer {
	//options := []grpctransport.ServerOption{
	//grpctransport.ServerErrorLogger(logger),
	//}
	return &grpcServer{
		// rello

		checklistwebhook: grpctransport.NewServer(
			ctx,
			endpoints.CheckListWebhookEndpoint,
			DecodeGRPCCheckListWebhookRequest,
			EncodeGRPCCheckListWebhookResponse,
			//append(options,grpctransport.ServerBefore(opentracing.FromGRPCRequest(tracer, "CheckListWebhook", logger)))...,
		),
	}
}

type grpcServer struct {
	checklistwebhook grpctransport.Handler
}

// Methods

func (s *grpcServer) CheckListWebhook(ctx context.Context, req *pb.ChecklistUpdate) (*pb.Empty, error) {
	_, rep, err := s.checklistwebhook.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.Empty), nil
}

// Server Decode

// DecodeGRPCCheckListWebhookRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC checklistwebhook request to a user-domain checklistwebhook request. Primarily useful in a server.
func DecodeGRPCCheckListWebhookRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ChecklistUpdate)
	return req, nil
}

// Client Decode

// DecodeGRPCCheckListWebhookResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC checklistwebhook reply to a user-domain checklistwebhook response. Primarily useful in a client.
func DecodeGRPCCheckListWebhookResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.Empty)
	return reply, nil
}

// Server Encode

// EncodeGRPCCheckListWebhookResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain checklistwebhook response to a gRPC checklistwebhook reply. Primarily useful in a server.
func EncodeGRPCCheckListWebhookResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.Empty)
	return resp, nil
}

// Client Encode

// EncodeGRPCCheckListWebhookRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain checklistwebhook request to a gRPC checklistwebhook request. Primarily useful in a client.
func EncodeGRPCCheckListWebhookRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ChecklistUpdate)
	return req, nil
}
