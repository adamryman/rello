package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	//"github.com/lightstep/lightstep-tracer-go"
	//stdopentracing "github.com/opentracing/opentracing-go"
	//zipkin "github.com/openzipkin/zipkin-go-opentracing"
	//appdashot "github.com/sourcegraph/appdash/opentracing"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	//"sourcegraph.com/sourcegraph/appdash"

	"github.com/pkg/errors"

	//"github.com/go-kit/kit/log"

	// This Service
	pb "github.com/adamryman/rello/rello-service"
	grpcclient "github.com/adamryman/rello/rello-service/generated/client/grpc"
	httpclient "github.com/adamryman/rello/rello-service/generated/client/http"
	clientHandler "github.com/adamryman/rello/rello-service/handlers/client"
	handler "github.com/adamryman/rello/rello-service/handlers/server"
)

var (
	_ = strconv.ParseInt
	_ = strings.Split
	_ = json.Compact
	_ = errors.Wrapf
	_ = pb.RegisterRelloServer
)

func main() {
	// The addcli presumes no service discovery system, and expects users to
	// provide the direct address of an addsvc. This presumption is reflected in
	// the addcli binary and the the client packages: the -transport.addr flags
	// and various client constructors both expect host:port strings. For an
	// example service with a client built on top of a service discovery system,
	// see profilesvc.

	var (
		httpAddr = flag.String("http.addr", "", "HTTP address of addsvc")
		grpcAddr = flag.String("grpc.addr", "", "gRPC (HTTP) address of addsvc")
		//zipkinAddr     = flag.String("zipkin.addr", "", "Enable Zipkin tracing via a Kafka Collector host:port")
		//appdashAddr    = flag.String("appdash.addr", "", "Enable Appdash tracing via an Appdash server host:port")
		//lightstepToken = flag.String("lightstep.token", "", "Enable LightStep tracing via a LightStep access token")
		method = flag.String("method", "checklistwebhook", "checklistwebhook")
	)

	var (
		flagModelCheckListWebhook  = flag.String("checklistwebhook.model", "", "")
		flagActionCheckListWebhook = flag.String("checklistwebhook.action", "", "")
	)
	flag.Parse()

	// This is a demonstration client, which supports multiple tracers.
	// Your clients will probably just use one tracer.
	//var tracer stdopentracing.Tracer
	//{
	//if *zipkinAddr != "" {
	//collector, err := zipkin.NewKafkaCollector(
	//strings.Split(*zipkinAddr, ","),
	//zipkin.KafkaLogger(log.NewNopLogger()),
	//)
	//if err != nil {
	//fmt.Fprintf(os.Stderr, "%v\n", err)
	//os.Exit(1)
	//}
	//tracer, err = zipkin.NewTracer(
	//zipkin.NewRecorder(collector, false, "localhost:8000", "addcli"),
	//)
	//if err != nil {
	//fmt.Fprintf(os.Stderr, "%v\n", err)
	//os.Exit(1)
	//}
	//} else if *appdashAddr != "" {
	//tracer = appdashot.NewTracer(appdash.NewRemoteCollector(*appdashAddr))
	//} else if *lightstepToken != "" {
	//tracer = lightstep.NewTracer(lightstep.Options{
	//AccessToken: *lightstepToken,
	//})
	//defer lightstep.FlushLightStepTracer(tracer)
	//} else {
	//tracer = stdopentracing.GlobalTracer() // no-op
	//}
	//}

	// This is a demonstration client, which supports multiple transports.
	// Your clients will probably just define and stick with 1 transport.

	var (
		service handler.Service
		err     error
	)
	if *httpAddr != "" {
		//service, err = httpclient.New(*httpAddr, tracer, log.NewNopLogger())
		service, err = httpclient.New(*httpAddr)
	} else if *grpcAddr != "" {
		conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while dialing grpc connection: %v", err)
			os.Exit(1)
		}
		defer conn.Close()
		service = grpcclient.New(conn /*, tracer, log.NewNopLogger()*/)
	} else {
		fmt.Fprintf(os.Stderr, "error: no remote address specified\n")
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	switch *method {

	case "checklistwebhook":

		var err error

		var ModelCheckListWebhook pb.Model
		if flagModelCheckListWebhook != nil && len(*flagModelCheckListWebhook) > 0 {
			err = json.Unmarshal([]byte(*flagModelCheckListWebhook), &ModelCheckListWebhook)
			if err != nil {
				panic(errors.Wrapf(err, "unmarshalling ModelCheckListWebhook from %v:", flagModelCheckListWebhook))
			}
		}

		var ActionCheckListWebhook pb.Action
		if flagActionCheckListWebhook != nil && len(*flagActionCheckListWebhook) > 0 {
			err = json.Unmarshal([]byte(*flagActionCheckListWebhook), &ActionCheckListWebhook)
			if err != nil {
				panic(errors.Wrapf(err, "unmarshalling ActionCheckListWebhook from %v:", flagActionCheckListWebhook))
			}
		}

		request, err := clientHandler.CheckListWebhook(ModelCheckListWebhook, ActionCheckListWebhook)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calling clientHandler.CheckListWebhook: %v\n", err)
			os.Exit(1)
		}

		v, err := service.CheckListWebhook(context.Background(), request)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calling service.CheckListWebhook: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Client Requested with:")
		fmt.Println(ModelCheckListWebhook, ActionCheckListWebhook)
		fmt.Println("Server Responded with:")
		fmt.Println(v)

	default:
		fmt.Fprintf(os.Stderr, "error: invalid method %q\n", method)
		os.Exit(1)
	}
}
