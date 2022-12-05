package main

import (
	"context"
	pb "externalscaler-license/externalscaler"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"net/http"
)

type ExternalScaler struct {
	Global bool
	Last   map[string]bool
	pb.UnsafeExternalScalerServer
}

func (e *ExternalScaler) IsActive(ctx context.Context, scaledObject *pb.ScaledObjectRef) (*pb.IsActiveResponse, error) {
	ns := scaledObject.Namespace

	if ns == "" {
		return nil, status.Error(codes.InvalidArgument, "Namespace must be specified")
	}

	isactive := e.isActiveTenant(ns)

	log.Printf("IsActive is %t for %s", isactive, ns)

	return &pb.IsActiveResponse{Result: isactive}, nil
}

// only used for external-push scalers
func (e *ExternalScaler) StreamIsActive(scaledObject *pb.ScaledObjectRef, epsServer pb.ExternalScaler_StreamIsActiveServer) error {
	return nil
}

func (e *ExternalScaler) GetMetricSpec(context.Context, *pb.ScaledObjectRef) (*pb.GetMetricSpecResponse, error) {

	return &pb.GetMetricSpecResponse{
		MetricSpecs: []*pb.MetricSpec{{
			MetricName: "licenseThreshold",
			TargetSize: 1,
		}},
	}, nil
}

func (e *ExternalScaler) GetMetrics(_ context.Context, metricRequest *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	return &pb.GetMetricsResponse{
		MetricValues: []*pb.MetricValue{{
			MetricName:  "licenseThreshold",
			MetricValue: 2,
		}},
	}, nil
}

// TODO: Get tenant info from license service
func (e *ExternalScaler) isActiveTenant(ns string) bool {
	isActive := false
	query := fmt.Sprintf("tenant=%s", ns)
	_, err := http.Get(fmt.Sprintf("http://licenses.default.svc.cluster.local:8080/tenant?%s", query))
	if err != nil {
		log.Printf("failed request %v", status.Error(codes.Internal, err.Error()))
		if val, ok := e.Last[ns]; ok {
			return val
		}
	}

	isActive = e.Global
	e.Last[ns] = isActive

	return isActive
}

// http scale toggle
func (e *ExternalScaler) SetGlobal(w http.ResponseWriter, r *http.Request) {
	log.Println("http called SetGlobal")
	e.Global = !e.Global
	io.WriteString(w, "Toggle Global \n")
}

func main() {
	grpcServer := grpc.NewServer()
	lis, _ := net.Listen("tcp", ":6000")
	es := &ExternalScaler{Global: true, Last: make(map[string]bool)}
	pb.RegisterExternalScalerServer(grpcServer, es)

	log.Println("listenting on :6000")
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	// http server
	http.HandleFunc("/", es.SetGlobal)
	log.Println("listenting on :3333")
	if err := http.ListenAndServe(":3333", nil); err != nil {
		log.Fatal(err)
	}

}
