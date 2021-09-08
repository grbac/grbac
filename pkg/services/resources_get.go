package services

import (
	"context"
	"encoding/base64"

	_ "embed"

	grbac "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"github.com/dgraph-io/dgo/v210"
	"github.com/grbac/grbac/pkg/graph"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AccessControlServerImpl) validateGetResource(ctx context.Context, txn *dgo.Txn, req *grbac.GetResourceRequest) error {
	if len(req.Name) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {resource name not defined}").Err()
	}

	// The resource name must be well formatted.
	if !isFullResourceName(req.Name) {
		return status.New(codes.InvalidArgument, "invalid argument {invalid resource name format}").Err()
	}

	return nil
}

// GetResource returns a resource.
func (s *AccessControlServerImpl) GetResource(ctx context.Context, req *grbac.GetResourceRequest) (*grbac.Resource, error) {
	txn := s.cli.NewReadOnlyTxn()
	if err := s.validateGetResource(ctx, txn, req); err != nil {
		return nil, err
	}

	// TODO(performance): GetResource should return only the resource name and parent (no policy).
	resp, err := graph.GetResource(ctx, txn, req.Name)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get resource [%s]", req.Name)
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	if resp == nil {
		return nil, status.New(codes.NotFound, "not found").Err()
	}

	resource := &grbac.Resource{
		Name: resp.Name,
	}

	resource.Etag, err = base64.StdEncoding.DecodeString(resp.ETag)
	if err != nil {
		logrus.WithError(err).Errorf("failed to decode resource etag [%s]", req.Name)
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	if resp.Parent != nil {
		resource.Parent = resp.Parent.Name
	}

	return resource, nil
}
