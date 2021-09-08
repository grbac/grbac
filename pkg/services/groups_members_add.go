package services

import (
	"context"

	grbac "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddGroupMember adds a member to a group.
func (s *AccessControlServerImpl) AddGroupMember(ctx context.Context, req *grbac.AddGroupMemberRequest) (*grbac.Group, error) {
	return nil, status.New(codes.Unimplemented, "unimplemented").Err()
}
