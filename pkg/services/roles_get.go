package services

import (
	"context"
	"encoding/base64"

	grbac "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"github.com/dgraph-io/dgo/v210"
	"github.com/grbac/grbac/pkg/graph"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AccessControlServerImpl) validateGetRole(ctx context.Context, txn *dgo.Txn, req *grbac.GetRoleRequest) error {
	if len(req.Name) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {role name not defined}").Err()
	}

	// The role name must be well formatted.
	if !isRole(req.Name) {
		return status.New(codes.InvalidArgument, "invalid argument {invalid role name format}").Err()
	}

	return nil
}

// GetRole returns a role.
func (s *AccessControlServerImpl) GetRole(ctx context.Context, req *grbac.GetRoleRequest) (*grbac.Role, error) {
	txn := s.cli.NewReadOnlyTxn()
	if err := s.validateGetRole(ctx, txn, req); err != nil {
		return nil, err
	}

	resp, err := graph.GetRole(ctx, txn, req.Name)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get role [%s]", req.Name)
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	if resp == nil {
		return nil, status.New(codes.NotFound, "not found").Err()
	}

	role := &grbac.Role{
		Name: resp.Name,
	}

	role.Etag, err = base64.StdEncoding.DecodeString(resp.ETag)
	if err != nil {
		logrus.WithError(err).Errorf("failed to decode role etag [%s]", req.Name)
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	for _, permission := range resp.Permissions {
		role.Permissions = append(role.Permissions, toPermissionId(permission.Name))
	}

	return role, nil
}
