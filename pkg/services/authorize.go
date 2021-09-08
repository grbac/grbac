package services

import (
	"context"
	"encoding/json"

	_ "embed"

	grbac "github.com/animeapis/go-genproto/grbac/v1alpha1"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO(christia-roggia): collapse into a single query as soon as dgraph
// allows `shortest` to be performed with multiple exit nodes.

//go:embed data/authorize.query.dql
var queryAuthorize string

func (s *AccessControlServerImpl) validateTestIamPolicy(ctx context.Context, req *grbac.TestIamPolicyRequest) error {
	if req.AccessTuple == nil {
		return status.New(codes.InvalidArgument, "invalid argument {access tuple not defined}").Err()
	}

	if len(req.AccessTuple.FullResourceName) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {full resource name not defined}").Err()
	}
	if len(req.AccessTuple.Permission) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {permission not defined}").Err()
	}
	if len(req.AccessTuple.Principal) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {principal not defined}").Err()
	}

	if !isUserMember(req.AccessTuple.Principal) && !isServiceAccountMember(req.AccessTuple.Principal) {
		return status.New(codes.InvalidArgument, "invalid argument {invalid principal name format}").Err()
	}

	return nil
}

// Checks whether a member has a specific permission for a specific resource.
// If not allowed an Unauthorized (403) error will be returned.
func (s *AccessControlServerImpl) TestIamPolicy(ctx context.Context, req *grbac.TestIamPolicyRequest) (*empty.Empty, error) {
	if err := s.validateTestIamPolicy(ctx, req); err != nil {
		return nil, err
	}

	m := map[string]string{
		"$resource":   req.AccessTuple.FullResourceName,
		"$permission": toPermissionName(req.AccessTuple.Permission),
	}

	if isUserMember(req.AccessTuple.Principal) {
		m["$principal"] = toUserName(req.AccessTuple.Principal)
	} else if isServiceAccountMember(req.AccessTuple.Principal) {
		m["$principal"] = toServiceAccountName(req.AccessTuple.Principal)
	}

	allUsers := map[string]string{
		"$principal":  allUsers,
		"$resource":   req.AccessTuple.FullResourceName,
		"$permission": toPermissionName(req.AccessTuple.Permission),
	}

	// Ask in parallel whether the user is allowed or allUsers is allowed.
	var isAllowed, isAllUsersAllowed bool
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		allowed, err := s.testIamPolicy(ctx, m)

		isAllowed = allowed
		return err
	})

	group.Go(func() error {
		allowed, err := s.testIamPolicy(ctx, allUsers)

		isAllUsersAllowed = allowed
		return err
	})

	if err := group.Wait(); err != nil {
		logrus.WithError(err).Errorf("failed to execute authorize query")
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	if isAllowed || isAllUsersAllowed {
		return &empty.Empty{}, nil
	}

	return nil, status.New(codes.PermissionDenied, "permission denied").Err()
}

func (s *AccessControlServerImpl) testIamPolicy(ctx context.Context, m map[string]string) (bool, error) {
	resp, err := s.cli.NewReadOnlyTxn().QueryWithVars(ctx, queryAuthorize, m)
	if err != nil {
		return false, err
	}

	payload := new(struct {
		Ok []json.RawMessage `json:"ok"`
	})

	if err := json.Unmarshal(resp.Json, &payload); err != nil {
		return false, err
	}

	if len(payload.Ok) == 0 {
		return false, nil
	}

	return true, nil
}
