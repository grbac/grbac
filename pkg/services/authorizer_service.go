package services

import (
	"context"
	"text/template"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"

	grbac "github.com/animeapis/go-genproto/grbac/v1alpha1"
)

type AccessControlServerConfig struct {
	DgraphHostname string
}

// NewAccessControlServer returns a new instance of AccessControl server.
func NewAccessControlServer(cfg *AccessControlServerConfig) (grbac.AccessControlServer, error) {
	connection, err := grpc.Dial(cfg.DgraphHostname, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &AccessControlServerImpl{
		cli:  dgo.NewDgraphClient(api.NewDgraphClient(connection)),
		conn: connection,
	}, nil
}

type AccessControlServerImpl struct {
	cli  *dgo.Dgraph
	conn *grpc.ClientConn
}

func (s *AccessControlServerImpl) Close() error {
	return s.conn.Close()
}

func (s *AccessControlServerImpl) delete(ctx context.Context, txn *dgo.Txn, queryTmpl, mutationTmpl *template.Template, data interface{}) error {
	query, err := ExecuteTemplate(queryTmpl, data)
	if err != nil {
		return err
	}

	mutation, err := ExecuteTemplate(mutationTmpl, data)
	if err != nil {
		return err
	}

	request := &api.Request{
		Query:     string(query),
		Mutations: []*api.Mutation{{DelNquads: mutation}},
		CommitNow: true,
	}

	if _, err := txn.Do(ctx, request); err != nil {
		return err
	}

	return nil
}

func (s *AccessControlServerImpl) create(ctx context.Context, txn *dgo.Txn, queryTmpl, mutationTmpl *template.Template, cond string, data interface{}) (*api.Response, error) {
	query, err := ExecuteTemplate(queryTmpl, data)
	if err != nil {
		return nil, err
	}

	mutation, err := ExecuteTemplate(mutationTmpl, data)
	if err != nil {
		return nil, err
	}

	request := &api.Request{
		Query:     string(query),
		Mutations: []*api.Mutation{{SetNquads: mutation, Cond: cond}},
		CommitNow: true,
	}

	return txn.Do(ctx, request)
}

func (s *AccessControlServerImpl) update(ctx context.Context, txn *dgo.Txn, queryTmpl, setTmpl, deleteTmpl *template.Template, data interface{}) error {
	query, err := ExecuteTemplate(queryTmpl, data)
	if err != nil {
		return err
	}

	setMutation, err := ExecuteTemplate(setTmpl, data)
	if err != nil {
		return err
	}

	deleteMutation, err := ExecuteTemplate(deleteTmpl, data)
	if err != nil {
		return err
	}

	request := &api.Request{
		Query:     string(query),
		Mutations: []*api.Mutation{{DelNquads: deleteMutation}, {SetNquads: setMutation}},
		CommitNow: true,
	}

	if _, err := txn.Do(ctx, request); err != nil {
		return err
	}

	return nil
}
