package bootstrap

import (
	"context"

	_ "embed"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"
)

//go:embed data/schema.rdf
var schema string

//go:embed data/system.all-users.query.rdf
var allUsersQuery string

//go:embed data/system.all-users.mutation.rdf
var allUsersMutation []byte

//go:embed data/system.all-users.condition.rdf
var allUsersCondition string

//go:embed data/system.animeshon.query.rdf
var animeshonQuery string

//go:embed data/system.animeshon.mutation.rdf
var animeshonMutation []byte

//go:embed data/system.animeshon.condition.rdf
var animeshonCondition string

func Schema(ctx context.Context, endpoint string) error {
	connection, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer connection.Close()

	op := &api.Operation{
		Schema: schema,
	}

	cli := dgo.NewDgraphClient(api.NewDgraphClient(connection))
	if err := cli.Alter(context.Background(), op); err != nil {
		return err
	}

	allUsers := &api.Request{
		Query: allUsersQuery,
		Mutations: []*api.Mutation{{
			Cond:      allUsersCondition,
			SetNquads: allUsersMutation,
		}},
		CommitNow: true,
	}

	if _, err := cli.NewTxn().Do(ctx, allUsers); err != nil {
		return err
	}

	animeshon := &api.Request{
		Query: animeshonQuery,
		Mutations: []*api.Mutation{{
			Cond:      animeshonCondition,
			SetNquads: animeshonMutation,
		}},
		CommitNow: true,
	}

	if _, err := cli.NewTxn().Do(ctx, animeshon); err != nil {
		return err
	}
	return nil
}
