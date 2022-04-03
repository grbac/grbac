package services

import (
	"context"
	"encoding/base64"
	"errors"
	"text/template"

	_ "embed"

	grbac "github.com/animeapis/go-genproto/grbac/v1alpha1"
	"github.com/dgraph-io/dgo/v210"
	"github.com/grbac/grbac/pkg/graph"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:embed data/groups/groups.create.query.go.tmpl
var queryCreateGroup string

//go:embed data/groups/groups.create.mutation.go.tmpl
var mutationCreateGroup string

var templateQueryCreateGroup = template.Must(
	template.New("QueryCreateGroup").Funcs(defaultFuncMap).Parse(queryCreateGroup),
)

var templateMutationCreateGroup = template.Must(
	template.New("MutationCreateGroup").Funcs(defaultFuncMap).Parse(mutationCreateGroup),
)

func (s *AccessControlServerImpl) validateCreateGroup(ctx context.Context, txn *dgo.Txn, req *grbac.CreateGroupRequest) error {
	// A group must be defined.
	if req.Group == nil {
		return status.New(codes.InvalidArgument, "invalid argument {group not defined}").Err()
	}

	// The group name must be defined.
	if len(req.Group.Name) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {group name not defined}").Err()
	}

	// The group name must be well formatted.
	if !isGroup(req.Group.Name) {
		return status.New(codes.InvalidArgument, "invalid argument {invalid group name format}").Err()
	}

	// The members must all exist and must have a valid type.
	for _, m := range req.Group.Members {
		memberFound, err := false, error(nil)
		if isGroupMember(m) {
			// TODO: should groups be allowed to include other groups?
			// TODO: if yes, a maximum path distance should be set to avoid too heavy queries.
			memberFound, err = graph.ExistsGroup(ctx, txn, toGroupName(m))
		} else if isUserMember(m) {
			memberFound, err = graph.ExistsSubject(ctx, txn, toUserName(m))
		} else if isServiceAccountMember(m) {
			memberFound, err = graph.ExistsSubject(ctx, txn, toServiceAccountName(m))
		} else if isAllUsersMember(m) {
			memberFound, err = graph.ExistsSubject(ctx, txn, allUsers)
		} else {
			return status.New(codes.InvalidArgument, "invalid argument {invalid member type}").Err()
		}

		if err != nil {
			logrus.WithError(err).Errorf("CreateGroup: failed to query group members")
			return status.New(codes.Internal, "internal error").Err()
		}

		if !memberFound {
			return status.New(codes.FailedPrecondition, "failed precondition {member does not exist}").Err()
		}
	}

	// The group must be new to avoid race conditions.
	groupFound, err := graph.ExistsGroup(ctx, txn, req.Group.Name)
	if err != nil {
		logrus.WithError(err).Errorf("CreateGroup: failed to query group")
		return status.New(codes.Internal, "internal error").Err()
	}

	if groupFound {
		return status.New(codes.AlreadyExists, "conflict").Err()
	}

	return nil
}

// CreateGroup creates a new group.
func (s *AccessControlServerImpl) CreateGroup(ctx context.Context, req *grbac.CreateGroupRequest) (*grbac.Group, error) {
	txn := s.cli.NewTxn()
	if err := s.validateCreateGroup(ctx, txn, req); err != nil {
		return nil, err
	}

	// TODO: etag should be generated according to the data structure.
	etag := []byte("TODO")

	data := struct {
		Group *grbac.Group
		ETag  string
	}{
		Group: req.GetGroup(),
		ETag:  base64.StdEncoding.EncodeToString(etag),
	}

	if _, err := s.create(ctx, txn, templateQueryCreateGroup, templateMutationCreateGroup, "", data); err != nil {
		if errors.Is(err, dgo.ErrAborted) {
			return nil, status.New(codes.Aborted, "transaction has been aborted").Err()
		}

		logrus.WithError(err).Errorf("CreateGroup: failed to execute dgraph call")
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	group := &grbac.Group{
		Name:    req.Group.Name,
		Members: req.Group.Members,
		Etag:    etag,
	}

	return group, nil
}
