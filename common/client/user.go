package client

import (
	"context"
	"fmt"

	"github.com/pixperk/notifly/common/auth"
	commonpb "github.com/pixperk/notifly/common/proto-gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type UserClient struct {
	conn      *grpc.ClientConn
	clientSvc commonpb.UserServiceClient
}

type AuthResp struct {
	Authenticated bool   `json:"authenticated"`
	Identifier    string `json:"identifier"`
}

type ValidateTokenResp struct {
	UserId     string `json:"id"`
	Identifier string `json:"identifier"`
}

func NewUserClient(url string) (*UserClient, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	clientSvc := commonpb.NewUserServiceClient(conn)
	return &UserClient{
		conn:      conn,
		clientSvc: clientSvc,
	}, nil
}

func (c *UserClient) Close() {
	c.conn.Close()
}

func (c *UserClient) SignUp(ctx context.Context, name, identifier, password string) (*AuthResp, error) {
	req := &commonpb.SignUpRequest{
		Name:       name,
		Identifier: identifier,
		Password:   password,
	}

	resp, err := c.clientSvc.SignUp(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("sign up failed: %w", err)
	}

	token := resp.Token
	// Create a new context with the auth token
	_ = metadata.AppendToOutgoingContext(ctx, auth.AuthMetadataKey, fmt.Sprintf("%s %s", auth.TokenPrefix, token))

	return &AuthResp{
		Authenticated: true,
		Identifier:    identifier,
	}, nil
}

func (c *UserClient) SignIn(ctx context.Context, identifier, password string) (*AuthResp, error) {
	req := &commonpb.SignInRequest{
		Identifier: identifier,
		Password:   password,
	}

	resp, err := c.clientSvc.SignIn(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("sign in failed: %w", err)
	}

	token := resp.Token
	// Create a new context with the auth token
	_ = metadata.AppendToOutgoingContext(ctx, auth.AuthMetadataKey, fmt.Sprintf("%s %s", auth.TokenPrefix, token))

	return &AuthResp{
		Authenticated: true,
		Identifier:    identifier,
	}, nil
}

func (c *UserClient) ValidateToken(ctx context.Context, token string) (*ValidateTokenResp, error) {
	req := &commonpb.ValidateTokenRequest{
		Token: token,
	}

	resp, err := c.clientSvc.ValidateToken(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	return &ValidateTokenResp{
		UserId:     resp.UserId,
		Identifier: resp.Identifier,
	}, nil
}
