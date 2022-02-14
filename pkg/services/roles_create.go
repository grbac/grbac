package services

import (
	"context"
	"encoding/base64"
	"errors"
	"text/template"

	_ "embed"

	grbac "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"github.com/dgraph-io/dgo/v210"
	"github.com/grbac/grbac/pkg/graph"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:embed data/roles/roles.create.query.go.tmpl
var queryCreateRole string

//go:embed data/roles/roles.create.mutation.go.tmpl
var mutationCreateRole string

var templateQueryCreateRole = template.Must(
	template.New("QueryCreateRole").Funcs(defaultFuncMap).Parse(queryCreateRole),
)

var templateMutationCreateRole = template.Must(
	template.New("MutationCreateRole").Funcs(defaultFuncMap).Parse(mutationCreateRole),
)

func (s *AccessControlServerImpl) validateCreateRole(ctx context.Context, txn *dgo.Txn, req *grbac.CreateRoleRequest) error {
	// A role must be defined.
	if req.Role == nil {
		return status.New(codes.InvalidArgument, "invalid argument {role not defined}").Err()
	}

	// The role name must be defined.
	if len(req.Role.Name) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {role name not defined}").Err()
	}

	// The role must include at least one permission.
	if len(req.Role.Permissions) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {role has no permissions}").Err()
	}

	// The permissions included in the role must exist.
	for _, permission := range req.Role.Permissions {
		permissionFound, err := graph.ExistsPermission(ctx, txn, toPermissionName(permission))
		if err != nil {
			logrus.WithError(err).Errorf("CreateRole: failed to query role permissions")
			return status.New(codes.Internal, "internal error").Err()
		}

		if !permissionFound {
			return status.New(codes.FailedPrecondition, "failed precondition {permission does not exist}").Err()
		}
	}

	// The role name must be well formatted.
	if !isRole(req.Role.Name) {
		return status.New(codes.InvalidArgument, "invalid argument {invalid role name format}").Err()
	}

	roleFound, err := graph.ExistsRole(ctx, txn, req.Role.Name)
	if err != nil {
		logrus.WithError(err).Errorf("failed to validate 'CreateRole' request")
		return status.New(codes.Internal, "internal error").Err()
	}

	if roleFound {
		return status.New(codes.AlreadyExists, "conflict").Err()
	}

	return nil
}

// CreateRole creates a new role.
func (s *AccessControlServerImpl) CreateRole(ctx context.Context, req *grbac.CreateRoleRequest) (*grbac.Role, error) {
	txn := s.cli.NewTxn()
	if err := s.validateCreateRole(ctx, txn, req); err != nil {
		return nil, err
	}

	// TODO: etag should be generated according to the data structure.
	etag := []byte("TODO")

	data := struct {
		Role *grbac.Role
		ETag string
	}{
		Role: req.GetRole(),
		ETag: base64.StdEncoding.EncodeToString(etag),
	}

	if err := s.create(ctx, txn, templateQueryCreateRole, templateMutationCreateRole, data); err != nil {
		if errors.Is(err, dgo.ErrAborted) {
			return nil, status.New(codes.Aborted, "transaction has been aborted").Err()
		}

		logrus.WithError(err).Errorf("CreateRole: failed to execute dgraph call")
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	role := &grbac.Role{
		Name:        req.Role.Name,
		Permissions: req.Role.Permissions,
		Etag:        etag,
	}

	return role, nil
}
