package services

import (
	"context"
	"encoding/base64"
	"errors"
	"text/template"

	_ "embed"

	"github.com/dgraph-io/dgo/v210"
	"github.com/grbac/grbac/pkg/graph"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:embed data/policies/policies.update.query.go.tmpl
var queryUpdatePolicy string

//go:embed data/policies/policies.update.set.go.tmpl
var setUpdatePolicy string

//go:embed data/policies/policies.update.delete.go.tmpl
var deleteUpdatePolicy string

var templateQueryUpdatePolicy = template.Must(
	template.New("QueryUpdatePolicy").Funcs(defaultFuncMap).Parse(queryUpdatePolicy),
)

var templateSetUpdatePolicy = template.Must(
	template.New("SetUpdatePolicy").Funcs(defaultFuncMap).Parse(setUpdatePolicy),
)

var templateDeleteUpdatePolicy = template.Must(
	template.New("DeleteUpdatePolicy").Funcs(defaultFuncMap).Parse(deleteUpdatePolicy),
)

func (s *AccessControlServerImpl) validateSetIamPolicy(ctx context.Context, txn *dgo.Txn, req *iam.SetIamPolicyRequest) error {
	// The resource name must be defined.
	if len(req.Resource) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {resource not defined}").Err()
	}

	// The full resource name must be well formatted.
	if !isFullResourceName(req.Resource) {
		return status.New(codes.InvalidArgument, "invalid argument {invalid resource name format}").Err()
	}

	// The resource policy is optional.
	if req.Policy == nil {
		return status.New(codes.InvalidArgument, "invalid argument {policy not defined}").Err()
	}

	// The policy version must be defined and valid.
	if req.Policy.Version != 1 {
		return status.New(codes.InvalidArgument, "invalid argument {invalid policy version}").Err()
	}

	roles := map[string]interface{}{}
	for _, i := range req.Policy.Bindings {
		// The binding role must be defined.
		if len(i.Role) == 0 {
			return status.New(codes.InvalidArgument, "invalid argument {role not defined}").Err()
		}

		// The role must exist.
		roleFound, err := graph.ExistsRole(ctx, txn, i.Role)
		if err != nil {
			logrus.WithError(err).Errorf("SetIamPolicy: failed to query role")
			return status.New(codes.Internal, "internal error").Err()
		}

		if !roleFound {
			return status.New(codes.InvalidArgument, "invalid argument {role does not exist}").Err()
		}

		// There must be at least one member in the binding.
		if len(i.Members) == 0 {
			return status.New(codes.InvalidArgument, "invalid argument {binding has no members}").Err()
		}

		// The same role should not appear multiple times across different bindings.
		if _, ok := roles[i.Role]; ok {
			return status.New(codes.InvalidArgument, "invalid argument {role defined multiple times}").Err()
		}
		roles[i.Role] = nil

		// The members must all exist and must have a known type.
		for _, m := range i.Members {
			memberFound := false
			if isGroupMember(m) {
				memberFound, err = graph.ExistsGroup(ctx, txn, toGroupName(m))
			} else if isUserMember(m) {
				memberFound, err = graph.ExistsSubject(ctx, txn, toUserName(m))
			} else if isServiceAccountMember(m) {
				memberFound, err = graph.ExistsSubject(ctx, txn, toServiceAccountName(m))
			} else if isAllUsersMember(m) {
				memberFound, err = graph.ExistsSubject(ctx, txn, allUsers)
			} else {
				return status.New(codes.InvalidArgument, "invalid argument {unknown member type}").Err()
			}

			if err != nil {
				logrus.WithError(err).Errorf("SetIamPolicy: failed to query binding members")
				return status.New(codes.Internal, "internal error").Err()
			}

			if !memberFound {
				return status.New(codes.InvalidArgument, "invalid argument {member does not exist}").Err()
			}
		}
	}

	// The resource must exist.
	resourceFound, err := graph.ExistsResource(ctx, txn, req.Resource)
	if err != nil {
		logrus.WithError(err).Errorf("SetIamPolicy: failed to query resource")
		return status.New(codes.Internal, "internal error").Err()
	}

	if !resourceFound {
		return status.New(codes.NotFound, "not found").Err()
	}

	return nil
}

// Sets the IAM policy that is attached to a generic resource.
func (s *AccessControlServerImpl) SetIamPolicy(ctx context.Context, req *iam.SetIamPolicyRequest) (*iam.Policy, error) {
	txn := s.cli.NewTxn()
	if err := s.validateSetIamPolicy(ctx, txn, req); err != nil {
		return nil, err
	}

	// TODO: etag should be generated according to the data structure.
	etag := []byte("TODO")

	data := struct {
		Resource string
		Policy   *iam.Policy
		ETag     string
	}{
		Resource: req.GetResource(),
		Policy:   req.GetPolicy(),
		ETag:     base64.StdEncoding.EncodeToString(etag),
	}

	if err := s.update(ctx, txn, templateQueryUpdatePolicy, templateSetUpdatePolicy, templateDeleteUpdatePolicy, data); err != nil {
		if errors.Is(err, dgo.ErrAborted) {
			return nil, status.New(codes.Aborted, "transaction has been aborted").Err()
		}

		logrus.WithError(err).Errorf("SetIamPolicy: failed to execute dgraph call")
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	policy := &iam.Policy{
		Version:  req.Policy.Version,
		Bindings: req.Policy.Bindings,
		Etag:     etag,
	}

	return policy, nil
}
