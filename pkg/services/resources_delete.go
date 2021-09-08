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

//go:embed data/resources/resources.delete.query.go.tmpl
var queryDeleteResource string

//go:embed data/resources/resources.delete.mutation.go.tmpl
var mutationDeleteResource string

var templateQueryDeleteResource = template.Must(
	template.New("QueryDeleteResource").Funcs(defaultFuncMap).Parse(queryDeleteResource),
)

var templateMutationDeleteResource = template.Must(
	template.New("MutationDeleteResource").Funcs(defaultFuncMap).Parse(mutationDeleteResource),
)

func (s *AccessControlServerImpl) validateDeleteResource(ctx context.Context, txn *dgo.Txn, req *grbac.DeleteResourceRequest) error {
	// The resource name must be defined.
	if len(req.Name) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {resource name not defined}").Err()
	}

	// The resource name must be well formatted.
	if !isFullResourceName(req.Name) {
		return status.New(codes.InvalidArgument, "invalid argument {invalid resource name format}").Err()
	}

	// The resource must exist.
	resourceFound, err := graph.ExistsResource(ctx, txn, req.Name)
	if err != nil {
		logrus.WithError(err).Errorf("DeleteResource: failed to query resource")
		return status.New(codes.Internal, "internal error").Err()
	}

	if !resourceFound {
		return status.New(codes.NotFound, "not found").Err()
	}

	// The resource must not have children before deletion.
	childrenFound, err := graph.HasChildren(ctx, txn, req.Name)
	if err != nil {
		logrus.WithError(err).Errorf("DeleteResource: failed to check if resource has children")
		return status.New(codes.Internal, "internal error").Err()
	}

	if childrenFound {
		return status.New(codes.FailedPrecondition, "failed precondition {resource has children}").Err()
	}

	return nil
}

// DeleteResource deletes a resource.
func (s *AccessControlServerImpl) DeleteResource(ctx context.Context, req *grbac.DeleteResourceRequest) (*empty.Empty, error) {
	txn := s.cli.NewTxn()
	if err := s.validateDeleteResource(ctx, txn, req); err != nil {
		return nil, err
	}

	data := struct {
		Name string
	}{
		Name: req.Name,
	}

	if err := s.delete(ctx, txn, templateQueryDeleteResource, templateMutationDeleteResource, data); err != nil {
		logrus.WithError(err).Errorf("DeleteResource: failed to execute dgraph call")
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	return &empty.Empty{}, nil
}
