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

//go:embed data/subjects/subjects.delete.query.go.tmpl
var queryDeleteSubject string

//go:embed data/subjects/subjects.delete.mutation.go.tmpl
var mutationDeleteSubject string

var templateQueryDeleteSubject = template.Must(
	template.New("QueryDeleteSubject").Funcs(defaultFuncMap).Parse(queryDeleteSubject),
)

var templateMutationDeleteSubject = template.Must(
	template.New("MutationDeleteSubject").Funcs(defaultFuncMap).Parse(mutationDeleteSubject),
)

func (s *AccessControlServerImpl) validateDeleteSubject(ctx context.Context, txn *dgo.Txn, req *grbac.DeleteSubjectRequest) error {
	// The subject name must be defined.
	if len(req.Name) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {subject name not defined}").Err()
	}

	// The subject name must be well formatted.
	if !isSubject(req.Name) {
		return status.New(codes.InvalidArgument, "invalid argument {invalid subject name format}").Err()
	}

	// The subject must exist.
	subjectFound, err := graph.ExistsSubject(ctx, txn, req.Name)
	if err != nil {
		logrus.WithError(err).Errorf("DeleteSubject: failed to query subject")
		return status.New(codes.Internal, "internal error").Err()
	}

	if !subjectFound {
		return status.New(codes.NotFound, "not found").Err()
	}

	return nil
}

// DeleteSubject deletes a subject.
func (s *AccessControlServerImpl) DeleteSubject(ctx context.Context, req *grbac.DeleteSubjectRequest) (*empty.Empty, error) {
	txn := s.cli.NewTxn()
	if err := s.validateDeleteSubject(ctx, txn, req); err != nil {
		return nil, err
	}

	data := struct {
		Name string
	}{
		Name: req.GetName(),
	}

	if err := s.delete(ctx, txn, templateQueryDeleteSubject, templateMutationDeleteSubject, data); err != nil {
		logrus.WithError(err).Errorf("DeleteSubject: failed to execute dgraph call")
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	return &empty.Empty{}, nil
}
