// Package grpc provides a gRPC client for the add service.
package grpc

import (
	//"time"

	//jujuratelimit "github.com/juju/ratelimit"
	//stdopentracing "github.com/opentracing/opentracing-go"
	//"github.com/sony/gobreaker"
	"google.golang.org/grpc"

	//"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	//"github.com/go-kit/kit/log"
	//"github.com/go-kit/kit/ratelimit"
	//"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	pb "github.com/adamryman/rello/rello-service"
	svc "github.com/adamryman/rello/rello-service/generated"
	handler "github.com/adamryman/rello/rello-service/handlers/server"
)

// New returns an AddService backed by a gRPC client connection. It is the
// responsibility of the caller to dial, and later close, the connection.
func New(conn *grpc.ClientConn /*, tracer stdopentracing.Tracer, logger log.Logger*/) handler.Service {
	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.

	//limiter := ratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(100, 100))

	var checklistwebhookEndpoint endpoint.Endpoint
	{
		checklistwebhookEndpoint = grpctransport.NewClient(
			conn,
			"rello.Rello",
			"CheckListWebhook",
			svc.EncodeGRPCCheckListWebhookRequest,
			svc.DecodeGRPCCheckListWebhookResponse,
			pb.Empty{},
			//grpctransport.ClientBefore(opentracing.FromGRPCRequest(tracer, "CheckListWebhook", logger)),
		).Endpoint()
		//checklistwebhookEndpoint = opentracing.TraceClient(tracer, "CheckListWebhook")(checklistwebhookEndpoint)
		//checklistwebhookEndpoint = limiter(checklistwebhookEndpoint)
		//checklistwebhookEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		//Name:    "CheckListWebhook",
		//Timeout: 30 * time.Second,
		//}))(checklistwebhookEndpoint)
	}

	return svc.Endpoints{

		CheckListWebhookEndpoint: checklistwebhookEndpoint,
	}
}
