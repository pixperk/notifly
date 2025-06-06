package trigger

import (
	"context"
	"fmt"
	"net"

	"github.com/pixperk/notifly/common"
	commonpb "github.com/pixperk/notifly/common/proto-gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
	commonpb.UnimplementedTriggerServiceServer
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	server := grpc.NewServer()
	commonpb.RegisterTriggerServiceServer(server, &grpcServer{service: s})
	reflection.Register(server)
	return server.Serve(lis)
}

func (s *grpcServer) Trigger(ctx context.Context, req *commonpb.NotificationRequest) (*commonpb.TriggerResponse, error) {
	event := common.NotificationEvent{
		Type: req.Type.Enum().String(),
	}
	resp, err := s.service.TriggerNotification(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to trigger notification: %w", err)
	}
	return resp, nil
}
