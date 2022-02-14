package services

import (
	"context"
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

//go:embed data/subjects/subjects.create.query.go.tmpl
var queryCreateSubject string

//go:embed data/subjects/subjects.create.mutation.go.tmpl
var mutationCreateSubject string

var templateQueryCreateSubject = template.Must(
	template.New("QueryCreateSubject").Funcs(defaultFuncMap).Parse(queryCreateSubject),
)

var templateMutationCreateSubject = template.Must(
	template.New("MutationCreateSubject").Funcs(defaultFuncMap).Parse(mutationCreateSubject),
)

func (s *AccessControlServerImpl) validateCreateSubject(ctx context.Context, txn *dgo.Txn, req *grbac.CreateSubjectRequest) error {
	// A subject must be defined.
	if req.Subject == nil {
		return status.New(codes.InvalidArgument, "invalid argument {subject not defined}").Err()
	}

	// The subject name must be defined.
	if len(req.Subject.Name) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {subject name not defined}").Err()
	}

	// The subject name must be well formatted.
	if !isSubject(req.Subject.Name) {
		return status.New(codes.InvalidArgument, "invalid argument {invalid subject name format}").Err()
	}

	// The subject must be new to avoid race conditions.
	subjectFound, err := graph.ExistsSubject(ctx, txn, req.Subject.Name)
	if err != nil {
		logrus.WithError(err).Errorf("failed to validate 'CreateSubject' request")
		return status.New(codes.Internal, "internal error").Err()
	}

	if subjectFound {
		return status.New(codes.AlreadyExists, "conflict").Err()
	}

	return nil
}

// CreateSubject creates a new subject.
func (s *AccessControlServerImpl) CreateSubject(ctx context.Context, req *grbac.CreateSubjectRequest) (*grbac.Subject, error) {
	txn := s.cli.NewTxn()
	if err := s.validateCreateSubject(ctx, txn, req); err != nil {
		return nil, err
	}

	data := struct {
		Subject *grbac.Subject
	}{
		Subject: req.GetSubject(),
	}

	if err := s.create(ctx, txn, templateQueryCreateSubject, templateMutationCreateSubject, data); err != nil {
		if errors.Is(err, dgo.ErrAborted) {
			return nil, status.New(codes.Aborted, "transaction has been aborted").Err()
		}

		logrus.WithError(err).Errorf("CreateSubject: failed to execute dgraph call")
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	return &grbac.Subject{Name: req.Subject.Name}, nil
}
