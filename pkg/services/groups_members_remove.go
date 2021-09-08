package services

import (
	"context"

	grbac "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RemoveGroupMember removes a member from a group.
func (s *AccessControlServerImpl) RemoveGroupMember(ctx context.Context, req *grbac.RemoveGroupMemberRequest) (*grbac.Group, error) {
	return nil, status.New(codes.Unimplemented, "unimplemented").Err()
}
