package util

import (
	"context"

	"github.com/pixperk/notifly/common/auth"
	"google.golang.org/grpc/metadata"
)

func WithToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, auth.AuthMetadataKey, token)
}
