package services

import (
	"context"
	"encoding/base64"

	"github.com/dgraph-io/dgo/v210"
	"github.com/grbac/grbac/pkg/graph"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AccessControlServerImpl) validateGetIamPolicy(ctx context.Context, txn *dgo.Txn, req *iam.GetIamPolicyRequest) error {
	if len(req.Resource) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {resource not defined}").Err()
	}

	// The full resource name must be well formatted.
	if !isFullResourceName(req.Resource) {
		return status.New(codes.InvalidArgument, "invalid argument {invalid resource name format}").Err()
	}

	return nil
}

// Gets the IAM policy that is attached to a generic resource.
func (s *AccessControlServerImpl) GetIamPolicy(ctx context.Context, req *iam.GetIamPolicyRequest) (*iam.Policy, error) {
	txn := s.cli.NewReadOnlyTxn()
	if err := s.validateGetIamPolicy(ctx, txn, req); err != nil {
		return nil, err
	}

	// TODO(performance): a new query should be used to query only the resource and its policy.
	resp, err := graph.GetResource(ctx, txn, req.GetResource())
	if err != nil {
		logrus.WithError(err).Errorf("failed to get resource [%s]", req.GetResource())
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	if resp == nil {
		return nil, status.New(codes.NotFound, "not found").Err()
	}

	if resp.Policy == nil {
		return &iam.Policy{}, nil
	}

	policy := &iam.Policy{
		Version: resp.Policy.Version,
	}

	policy.Etag, err = base64.StdEncoding.DecodeString(resp.Policy.ETag)
	if err != nil {
		logrus.WithError(err).Errorf("failed to decode policy etag [%s]", req.Resource)
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	for _, i := range resp.Policy.Bindings {
		if i.Role == nil {
			logrus.Warningf("found binding with no role in resource [%s]", resp.Name)
			continue
		}

		binding := &iam.Binding{
			Role: i.Role.Name,
		}

		binding.Members, err = members(i.Members)
		if err != nil {
			logrus.WithError(err).Errorf("failed to get binding members [%s:%s]", req.Resource, i.Role.Name)
			return nil, status.New(codes.Internal, "internal error").Err()
		}

		policy.Bindings = append(policy.Bindings, binding)
	}

	return policy, nil
}
