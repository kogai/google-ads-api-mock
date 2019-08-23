package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/genproto/googleapis/ads/googleads/v1/services"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) Search(ctx context.Context, in *services.SearchGoogleAdsRequest) (*services.SearchGoogleAdsResponse, error) {
	log.Printf("Received(Service): %v", in)
	return &services.SearchGoogleAdsResponse{}, nil
}

func (s *server) Mutate(ctx context.Context, in *services.MutateGoogleAdsRequest) (*services.MutateGoogleAdsResponse, error) {
	log.Printf("Received(Mutate): %v", in)
	return &services.MutateGoogleAdsResponse{}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	services.RegisterGoogleAdsServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
