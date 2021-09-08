package services

import (
	"context"

	grbac "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TransferResource transfers a resource to a new parent.
func (s *AccessControlServerImpl) TransferResource(ctx context.Context, req *grbac.TransferResourceRequest) (*grbac.Resource, error) {
	return nil, status.New(codes.Unimplemented, "unimplemented").Err()
}
