package graphql

import "context"

type queryResolver struct {
	server *GraphQLServer
}

func (r *queryResolver) HealthCheck(ctx context.Context) (string, error) {
	return "OK", nil
}
