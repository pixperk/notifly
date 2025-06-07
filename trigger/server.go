package trigger

import (
	"context"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/pixperk/notifly/common"
	"github.com/pixperk/notifly/common/auth"
	commonpb "github.com/pixperk/notifly/common/proto-gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
	commonpb.UnimplementedTriggerServiceServer
}

func ListenGRPC(s Service, port int, tokenMaker auth.TokenMaker) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	server := grpc.NewServer(
		grpc.UnaryInterceptor(auth.AuthUnaryInterceptor(tokenMaker)))
	commonpb.RegisterTriggerServiceServer(server, &grpcServer{service: s})
	reflection.Register(server)
	return server.Serve(lis)
}

func (s *grpcServer) TriggerNotification(ctx context.Context, req *commonpb.NotificationRequest) (*commonpb.TriggerResponse, error) {

	authPayload, err := auth.GetAuthPayload(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get auth payload: %w", err)
	}

	notificationId := uuid.New()

	event := common.NotificationEvent{
		NotificationId: notificationId,
		Type:           req.Type.Enum().String(),
		Recipient:      req.Recipient,
		Subject:        req.Subject,
		Body:           req.Body,
		TriggerBy:      authPayload.Identifier,
	}
	err = s.service.TriggerNotification(event)
	if err != nil {
		return nil, fmt.Errorf("failed to trigger notification: %w", err)
	}
	return &commonpb.TriggerResponse{
		Status:         commonpb.TriggerResponse_QUEUED,
		Message:        "Notification being sent",
		TriggerBy:      authPayload.Identifier,
		NotificationId: notificationId.String(),
	}, nil
}
