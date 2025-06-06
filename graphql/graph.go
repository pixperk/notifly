package graphql

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/pixperk/notifly/common/client"
	"github.com/pixperk/notifly/graphql/generated"
)

type GraphQLServer struct {
	userClient    *client.UserClient
	triggerClient *client.TriggerClient
	token         string
}

func NewGraphQLServer(userUrl, accountUrl string) (*GraphQLServer, error) {
	userClient, err := client.NewUserClient(userUrl)
	if err != nil {
		return nil, err
	}
	triggerClient, err := client.NewTriggerClient(accountUrl)
	if err != nil {
		userClient.Close()
		return nil, err
	}

	return &GraphQLServer{
		userClient:    userClient,
		triggerClient: triggerClient,
	}, nil
}

func (s *GraphQLServer) Mutation() generated.MutationResolver {
	return &mutationResolver{server: s}
}

func (s *GraphQLServer) Query() generated.QueryResolver {
	return &queryResolver{server: s}
}

func (s *GraphQLServer) ToExecutableSchema() graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: s,
	})
}
