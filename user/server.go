package user

import (
	"context"
	"fmt"
	"net"

	commonpb "github.com/pixperk/notifly/common/proto-gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
	commonpb.UnimplementedUserServiceServer
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	server := grpc.NewServer()
	commonpb.RegisterUserServiceServer(server, &grpcServer{service: s})
	reflection.Register(server)
	return server.Serve(lis)
}

func (s *grpcServer) SignUp(ctx context.Context, req *commonpb.SignUpRequest) (*commonpb.AuthResponse, error) {
	authResponse, err := s.service.SignUp(ctx, req.Name, req.Identifier, req.Password)
	if err != nil {
		return nil, err
	}
	return &commonpb.AuthResponse{
		Identifier: authResponse.Identifier,
		Token:      authResponse.Token,
	}, nil
}

func (s *grpcServer) SignIn(ctx context.Context, req *commonpb.SignInRequest) (*commonpb.AuthResponse, error) {
	authResponse, err := s.service.SignIn(ctx, req.Identifier, req.Password)
	if err != nil {
		return nil, err
	}
	return &commonpb.AuthResponse{
		Identifier: authResponse.Identifier,
		Token:      authResponse.Token,
	}, nil
}

func (s *grpcServer) ValidateToken(ctx context.Context, req *commonpb.ValidateTokenRequest) (*commonpb.ValidateTokenResponse, error) {
	validateResponse, err := s.service.ValidateToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &commonpb.ValidateTokenResponse{
		Identifier: validateResponse.Identifier,
		UserId:     validateResponse.UserID.String(),
	}, nil
}
