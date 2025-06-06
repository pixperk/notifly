package graphql

import (
	"context"

	"github.com/pixperk/notifly/common"
	"github.com/pixperk/notifly/graphql/models"
	"github.com/pixperk/notifly/graphql/util"
)

type mutationResolver struct {
	server *GraphQLServer
}

func (r *mutationResolver) SignUp(ctx context.Context, input models.SignUpInput) (*models.AuthResp, error) {
	authResp, err := r.server.userClient.SignUp(ctx, input.Name, input.Identifier, input.Password)
	if err != nil {
		return nil, err
	}

	r.server.token = authResp.Token

	return &models.AuthResp{
		Authenticated: true,
		Identifier:    authResp.Identifier,
	}, nil
}

func (r *mutationResolver) SignIn(ctx context.Context, input models.SignInInput) (*models.AuthResp, error) {
	authResp, err := r.server.userClient.SignIn(ctx, input.Identifier, input.Password)
	if err != nil {
		return nil, err
	}

	r.server.token = authResp.Token

	return &models.AuthResp{
		Authenticated: true,
		Identifier:    authResp.Identifier,
	}, nil
}

func (r *mutationResolver) ValidateToken(ctx context.Context, input models.ValidateTokenInput) (*models.ValidateTokenResp, error) {
	resp, err := r.server.userClient.ValidateToken(ctx, input.Token)
	if err != nil {
		return nil, err
	}

	return &models.ValidateTokenResp{
		UserID:     resp.UserId,
		Identifier: resp.Identifier,
	}, nil
}

func (r *mutationResolver) TriggerNotification(ctx context.Context, input models.NotificationInput) (*models.TriggerNotificationResp, error) {
	var intputSubject string
	if input.Subject == nil {
		intputSubject = "Default Subject"
	} else {
		intputSubject = *input.Subject
	}

	event := common.NotificationEvent{
		Type:      input.Type.String(),
		Recipient: input.Recipient,
		Subject:   intputSubject,
		Body:      input.Body,
	}

	ctxWithToken := util.WithToken(ctx, r.server.token)

	resp, err := r.server.triggerClient.TriggerNotification(ctxWithToken, event)
	if err != nil {
		return nil, err
	}
	return &models.TriggerNotificationResp{
		NotificationID: resp.NotificationId,
		Status:         resp.Status,
		Message:        resp.Message,
		TriggerBy:      resp.TriggerBy,
	}, nil
}
