package main

import (
	"context"
	"log"
	"net"
	"os"
	"strings"

	"google.golang.org/genproto/googleapis/ads/googleads/v1/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	defaultPort = ":50051"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
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

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// Perform the token validation here. For the sake of this example, the code
	// here forgoes any of the usual OAuth2 token validation and instead checks
	// for a token matching an arbitrary string.
	return token == "some-secret-token"
}

func ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	} else {
		port = ":" + port
	}
	lis, err := net.Listen("tcp", ":443")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	crd, err := credentials.NewServerTLSFromFile("./localhost+2.pem", "./localhost+2-key.pem")
	if err != nil {
		log.Fatalf("failed to lookup credeintials: %v", err)
	}
	log.Printf("%v", crd)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(ensureValidToken),
		grpc.Creds(crd))
	services.RegisterGoogleAdsServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
