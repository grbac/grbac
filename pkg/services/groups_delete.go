package services

import (
	"context"
	"text/template"

	_ "embed"

	grbac "github.com/animeapis/go-genproto/grbac/v1alpha1"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/dgraph-io/dgo/v210"
	"github.com/grbac/grbac/pkg/graph"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:embed data/groups/groups.delete.query.go.tmpl
var queryDeleteGroup string

//go:embed data/groups/groups.delete.mutation.go.tmpl
var mutationDeleteGroup string

var templateQueryDeleteGroup = template.Must(
	template.New("QueryDeleteGroup").Funcs(defaultFuncMap).Parse(queryDeleteGroup),
)

var templateMutationDeleteGroup = template.Must(
	template.New("MutationDeleteGroup").Funcs(defaultFuncMap).Parse(mutationDeleteGroup),
)

func (s *AccessControlServerImpl) validateDeleteGroup(ctx context.Context, txn *dgo.Txn, req *grbac.DeleteGroupRequest) error {
	// The group name must be defined.
	if len(req.Name) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {group name not defined}").Err()
	}

	// The group name must be well formatted.
	if !isGroup(req.Name) {
		return status.New(codes.InvalidArgument, "invalid argument {invalid group name format}").Err()
	}

	// The group must exist.
	groupFound, err := graph.ExistsGroup(ctx, txn, req.Name)
	if err != nil {
		logrus.WithError(err).Errorf("DeleteGroup: failed to query group")
		return status.New(codes.Internal, "internal error").Err()
	}

	if !groupFound {
		return status.New(codes.NotFound, "not found").Err()
	}

	return nil
}

// DeleteGroup deletes a group.
func (s *AccessControlServerImpl) DeleteGroup(ctx context.Context, req *grbac.DeleteGroupRequest) (*empty.Empty, error) {
	txn := s.cli.NewTxn()
	if err := s.validateDeleteGroup(ctx, txn, req); err != nil {
		return nil, err
	}

	data := struct {
		Name string
	}{
		Name: req.GetName(),
	}

	if err := s.delete(ctx, txn, templateQueryDeleteGroup, templateMutationDeleteGroup, data); err != nil {
		logrus.WithError(err).Errorf("DeleteGroup: failed to execute dgraph call")
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	return &empty.Empty{}, nil
}
