package services

import (
	"context"
	"encoding/base64"
	"text/template"

	_ "embed"

	grbac "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"github.com/dgraph-io/dgo/v210"
	"github.com/grbac/grbac/pkg/fieldmask"
	"github.com/grbac/grbac/pkg/graph"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:embed data/roles/roles.update.query.go.tmpl
var queryUpdateRole string

//go:embed data/roles/roles.update.set.go.tmpl
var setUpdateRole string

//go:embed data/roles/roles.update.delete.go.tmpl
var deleteUpdateRole string

var templateQueryUpdateRole = template.Must(
	template.New("QueryUpdateRole").Funcs(defaultFuncMap).Parse(queryUpdateRole),
)

var templateSetUpdateRole = template.Must(
	template.New("SetUpdateRole").Funcs(defaultFuncMap).Parse(setUpdateRole),
)

var templateDeleteUpdateRole = template.Must(
	template.New("DeleteUpdateRole").Funcs(defaultFuncMap).Parse(deleteUpdateRole),
)

func (s *AccessControlServerImpl) validateUpdateRole(ctx context.Context, txn *dgo.Txn, req *grbac.UpdateRoleRequest) error {
	// A role must be defined.
	if req.Role == nil {
		return status.New(codes.InvalidArgument, "invalid argument {role not defined}").Err()
	}

	// The role name must be defined.
	if len(req.Role.Name) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {role name not defined}").Err()
	}

	// The role name must be well formatted.
	if !isRole(req.Role.Name) {
		return status.New(codes.InvalidArgument, "invalid argument {invalid role name format}").Err()
	}

	// The update field mask must contain valid paths.
	for _, path := range req.GetUpdateMask().GetPaths() {
		switch path {
		case "role", "role.permissions":
		default:
			return status.New(codes.InvalidArgument, "invalid argument {invalid field mask}").Err()
		}
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

	// The role must exist.
	roleFound, err := graph.ExistsRole(ctx, txn, req.Role.Name)
	if err != nil {
		logrus.WithError(err).Errorf("UpdateRole: failed to query role")
		return status.New(codes.Internal, "internal error").Err()
	}

	if !roleFound {
		return status.New(codes.NotFound, "not found").Err()
	}

	return nil
}

// UpdateRole updates a role with a field mask.
func (s *AccessControlServerImpl) UpdateRole(ctx context.Context, req *grbac.UpdateRoleRequest) (*grbac.Role, error) {
	txn := s.cli.NewTxn()
	if err := s.validateUpdateRole(ctx, txn, req); err != nil {
		return nil, err
	}

	// TODO: etag should be generated according to the data structure.
	etag := []byte("TODO")

	fieldmask := fieldmask.NewFieldMask(req.GetUpdateMask())

	data := struct {
		Role      *grbac.Role
		FieldMask func(string) bool
		ETag      string
	}{
		Role:      req.GetRole(),
		FieldMask: fieldmask.Contains,
		ETag:      base64.StdEncoding.EncodeToString(etag),
	}

	if err := s.update(ctx, txn, templateQueryUpdateRole, templateSetUpdateRole, templateDeleteUpdateRole, data); err != nil {
		logrus.WithError(err).Errorf("UpdateRole: failed to execute dgraph call")
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	// TODO: merge missing fields (not included in the update mask) with the role in dgraph.
	role := &grbac.Role{
		Name:        req.Role.Name,
		Permissions: req.Role.Permissions,
		Etag:        etag,
	}

	return role, nil
}
