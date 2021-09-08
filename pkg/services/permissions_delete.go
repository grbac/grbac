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

//go:embed data/permissions/permissions.delete.query.go.tmpl
var queryDeletePermission string

//go:embed data/permissions/permissions.delete.mutation.go.tmpl
var mutationDeletePermission string

var templateQueryDeletePermission = template.Must(
	template.New("QueryDeletePermission").Funcs(defaultFuncMap).Parse(queryDeletePermission),
)

var templateMutationDeletePermission = template.Must(
	template.New("MutationDeletePermission").Funcs(defaultFuncMap).Parse(mutationDeletePermission),
)

func (s *AccessControlServerImpl) validateDeletePermission(ctx context.Context, txn *dgo.Txn, req *grbac.DeletePermissionRequest) error {
	// The permission name must be defined.
	if len(req.Name) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {permission name not defined}").Err()
	}

	// The permission name must be well formatted.
	if !isPermission(req.Name) {
		return status.New(codes.InvalidArgument, "invalid argument {invalid permission name format}").Err()
	}

	// The permission must exist.
	permissionFound, err := graph.ExistsPermission(ctx, txn, req.Name)
	if err != nil {
		logrus.WithError(err).Errorf("DeletePermission: failed to query permission")
		return status.New(codes.Internal, "internal error").Err()
	}

	if !permissionFound {
		return status.New(codes.NotFound, "not found").Err()
	}

	return nil
}

// DeletePermission deletes a permission.
func (s *AccessControlServerImpl) DeletePermission(ctx context.Context, req *grbac.DeletePermissionRequest) (*empty.Empty, error) {
	txn := s.cli.NewTxn()
	if err := s.validateDeletePermission(ctx, txn, req); err != nil {
		return nil, err
	}

	data := struct {
		Name string
	}{
		Name: req.GetName(),
	}

	if err := s.delete(ctx, txn, templateQueryDeletePermission, templateMutationDeletePermission, data); err != nil {
		logrus.WithError(err).Errorf("DeletePermission: failed to execute dgraph call")
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	return &empty.Empty{}, nil
}
