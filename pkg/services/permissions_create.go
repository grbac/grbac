package services

import (
	"context"
	"text/template"

	_ "embed"

	grbac "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"github.com/dgraph-io/dgo/v210"
	"github.com/grbac/grbac/pkg/graph"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:embed data/permissions/permissions.create.query.go.tmpl
var queryCreatePermission string

//go:embed data/permissions/permissions.create.mutation.go.tmpl
var mutationCreatePermission string

var templateQueryCreatePermission = template.Must(
	template.New("QueryCreatePermission").Funcs(defaultFuncMap).Parse(queryCreatePermission),
)

var templateMutationCreatePermission = template.Must(
	template.New("MutationCreatePermission").Funcs(defaultFuncMap).Parse(mutationCreatePermission),
)

func (s *AccessControlServerImpl) validateCreatePermission(ctx context.Context, txn *dgo.Txn, req *grbac.CreatePermissionRequest) error {
	// A permission must be defined.
	if req.Permission == nil {
		return status.New(codes.InvalidArgument, "invalid argument {permission not defined}").Err()
	}

	// The permission name must be defined.
	if len(req.Permission.Name) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {permission name not defined}").Err()
	}

	// The permission name must be well formatted.
	if !isPermission(req.Permission.Name) || !isValidPermissionId(req.Permission.Name) {
		return status.New(codes.InvalidArgument, "invalid argument {invalid permission name format}").Err()
	}

	// The permission must be new to avoid race conditions.
	permissionFound, err := graph.ExistsPermission(ctx, txn, req.Permission.Name)
	if err != nil {
		logrus.WithError(err).Errorf("failed to validate 'CreatePermission' request")
		return status.New(codes.Internal, "internal error").Err()
	}

	if permissionFound {
		return status.New(codes.AlreadyExists, "conflict").Err()
	}

	return nil
}

// CreatePermission creates a new permission.
func (s *AccessControlServerImpl) CreatePermission(ctx context.Context, req *grbac.CreatePermissionRequest) (*grbac.Permission, error) {
	txn := s.cli.NewTxn()
	if err := s.validateCreatePermission(ctx, txn, req); err != nil {
		return nil, err
	}

	data := struct {
		Permission *grbac.Permission
	}{
		Permission: req.GetPermission(),
	}

	if err := s.create(ctx, txn, templateQueryCreatePermission, templateMutationCreatePermission, data); err != nil {
		logrus.WithError(err).Errorf("CreatePermission: failed to execute dgraph call")
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	return &grbac.Permission{Name: req.Permission.Name}, nil
}
