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

func (s *AccessControlServerImpl) validateGetGroup(ctx context.Context, txn *dgo.Txn, req *grbac.GetGroupRequest) error {
	if len(req.Name) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {group name not defined}").Err()
	}

	// The group name must be well formatted.
	if !isGroup(req.Name) {
		return status.New(codes.InvalidArgument, "invalid argument {invalid group name format}").Err()
	}

	return nil
}

// GetGroup returns a group.
func (s *AccessControlServerImpl) GetGroup(ctx context.Context, req *grbac.GetGroupRequest) (*grbac.Group, error) {
	txn := s.cli.NewReadOnlyTxn()
	if err := s.validateGetGroup(ctx, txn, req); err != nil {
		return nil, err
	}

	resp, err := graph.GetGroup(ctx, txn, req.GetName())
	if err != nil {
		logrus.WithError(err).Errorf("failed to get group [%s]", req.GetName())
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	if resp == nil {
		return nil, status.New(codes.NotFound, "not found").Err()
	}

	group := &grbac.Group{
		Name: resp.Name,
	}

	group.Etag, err = base64.StdEncoding.DecodeString(resp.ETag)
	if err != nil {
		logrus.WithError(err).Errorf("failed to decode resource etag [%s]", req.Name)
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	group.Members, err = members(resp.Members)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get group members [%s]", req.Name)
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	return group, nil
}
