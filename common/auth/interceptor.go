package auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	AuthMetadataKey = "authorization"
)

type contextKey string

const authPayloadKey contextKey = "auth_payload"

func AuthUnaryInterceptor(tokenMaker TokenMaker) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		tokenArray := md[AuthMetadataKey]

		if len(tokenArray) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing auth token")
		}

		token := tokenArray[0]

		payload, err := tokenMaker.VerifyToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
		}

		newCtx := context.WithValue(ctx, authPayloadKey, payload)
		return handler(newCtx, req)

	}
}

func GetAuthPayload(ctx context.Context) (*Payload, error) {
	payload, ok := ctx.Value(authPayloadKey).(*Payload)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "auth payload not found in context")
	}
	return payload, nil
}
