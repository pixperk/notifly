package client

import (
	"context"
	"fmt"

	"github.com/pixperk/notifly/common"
	commonpb "github.com/pixperk/notifly/common/proto-gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TriggerClient struct {
	conn      *grpc.ClientConn
	clientSvc commonpb.TriggerServiceClient
}

type TriggerNotificationResp struct {
	NotificationId string `json:"notification_id"`
	Status         string `json:"status"`
	Message        string `json:"message"`
	TriggerBy      string `json:"triggered_by"`
}

func NewTriggerClient(url string) (*TriggerClient, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	clientSvc := commonpb.NewTriggerServiceClient(conn)
	return &TriggerClient{
		conn:      conn,
		clientSvc: clientSvc,
	}, nil
}

func (c *TriggerClient) Close() {
	c.conn.Close()
}

func (c *TriggerClient) TriggerNotification(ctx context.Context, event common.NotificationEvent) (*TriggerNotificationResp, error) {

	var notifType commonpb.NotificationRequest_NotificationType

	if val, ok := commonpb.NotificationRequest_NotificationType_value[event.Type]; ok {
		notifType = commonpb.NotificationRequest_NotificationType(val)
	} else {
		return nil, fmt.Errorf("invalid notification type: %s", event.Type)
	}

	req := &commonpb.NotificationRequest{
		Type:      notifType,
		Recipient: event.Recipient,
		Subject:   event.Subject,
		Body:      event.Body,
	}

	resp, err := c.clientSvc.TriggerNotification(ctx, req)

	if err != nil {
		return nil, err
	}

	return &TriggerNotificationResp{
		NotificationId: resp.NotificationId,
		Status:         resp.Status.String(),
		Message:        resp.Message,
		TriggerBy:      resp.TriggerBy,
	}, nil
}
