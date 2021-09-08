package services

import (
	"context"
	"text/template"

	_ "embed"

	grbac "github.com/animeapis/go-genproto/grbac/v1alpha1"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/dgraph-io/dgo/v210"
	"github.com/grbac/grbac/pkg/graph"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:embed data/roles/roles.delete.query.go.tmpl
var queryDeleteRole string

//go:embed data/roles/roles.delete.mutation.go.tmpl
var mutationDeleteRole string

var templateQueryDeleteRole = template.Must(
	template.New("QueryDeleteRole").Funcs(defaultFuncMap).Parse(queryDeleteRole),
)

var templateMutationDeleteRole = template.Must(
	template.New("MutationDeleteRole").Funcs(defaultFuncMap).Parse(mutationDeleteRole),
)

func (s *AccessControlServerImpl) validateDeleteRole(ctx context.Context, txn *dgo.Txn, req *grbac.DeleteRoleRequest) error {
	// The role name must be defined.
	if len(req.Name) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {role name not defined}").Err()
	}

	// The role name must be well formatted.
	if !isRole(req.Name) {
		return status.New(codes.InvalidArgument, "invalid argument {invalid role name format}").Err()
	}

	// The role must exist.
	roleFound, err := graph.ExistsRole(ctx, txn, req.Name)
	if err != nil {
		logrus.WithError(err).Errorf("DeleteRole: failed to query role")
		return status.New(codes.Internal, "internal error").Err()
	}

	if !roleFound {
		return status.New(codes.NotFound, "not found").Err()
	}

	return nil
}

// DeleteRole deletes a role.
func (s *AccessControlServerImpl) DeleteRole(ctx context.Context, req *grbac.DeleteRoleRequest) (*empty.Empty, error) {
	txn := s.cli.NewTxn()
	if err := s.validateDeleteRole(ctx, txn, req); err != nil {
		return nil, err
	}

	data := struct {
		Name string
	}{
		Name: req.GetName(),
	}

	if err := s.delete(ctx, txn, templateQueryDeleteRole, templateMutationDeleteRole, data); err != nil {
		logrus.WithError(err).Errorf("DeleteRole: failed to execute dgraph call")
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	return &empty.Empty{}, nil
}
