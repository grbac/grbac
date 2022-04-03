package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"text/template"

	_ "embed"

	grbac "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"github.com/dgraph-io/dgo/v210"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:embed data/resources/resources.create.query.go.tmpl
var queryCreateResource string

//go:embed data/resources/resources.create.mutation.go.tmpl
var mutationCreateResource string

var conditionCreateResource = "@if(not(eq(len(parent), 0)) AND eq(len(resource), 0))"

var templateQueryCreateResource = template.Must(
	template.New("QueryCreateResource").Funcs(defaultFuncMap).Parse(queryCreateResource),
)

var templateMutationCreateResource = template.Must(
	template.New("MutationCreateResource").Funcs(defaultFuncMap).Parse(mutationCreateResource),
)

func (s *AccessControlServerImpl) validateCreateResource(ctx context.Context, txn *dgo.Txn, req *grbac.CreateResourceRequest) error {
	// A resource must be defined.
	if req.Resource == nil {
		return status.New(codes.InvalidArgument, "invalid argument {resource not defined}").Err()
	}

	// The resource name must be defined.
	if len(req.Resource.Name) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {resource name not defined}").Err()
	}

	// The resource name must be well formatted.
	if !isFullResourceName(req.Resource.Name) {
		return status.New(codes.InvalidArgument, "invalid argument {invalid resource name format}").Err()
	}

	// The parent name must be defined.
	if len(req.Resource.Parent) == 0 {
		return status.New(codes.InvalidArgument, "invalid argument {parent name not defined}").Err()
	}

	// The parent name must be well formatted.
	if !isFullResourceName(req.Resource.Parent) {
		return status.New(codes.InvalidArgument, "invalid argument {invalid parent name format}").Err()
	}

	return nil
}

// CreateResource creates a new resource.
func (s *AccessControlServerImpl) CreateResource(ctx context.Context, req *grbac.CreateResourceRequest) (*grbac.Resource, error) {
	txn := s.cli.NewTxn()
	if err := s.validateCreateResource(ctx, txn, req); err != nil {
		return nil, err
	}

	// TODO: etag should be generated according to the data structure.
	etag := []byte("TODO")

	data := struct {
		Resource *grbac.Resource
		ETag     string
	}{
		Resource: req.GetResource(),
		ETag:     base64.StdEncoding.EncodeToString(etag),
	}

	resp, err := s.create(ctx, txn, templateQueryCreateResource, templateMutationCreateResource, conditionCreateResource, data)
	if err != nil {
		if errors.Is(err, dgo.ErrAborted) {
			return nil, status.New(codes.Aborted, "transaction has been aborted").Err()
		}

		logrus.WithError(err).Errorf("CreateResource: failed to execute dgraph call")
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	var m map[string][]map[string]string
	if err := json.Unmarshal(resp.Json, &m); err != nil {
		logrus.WithError(err).Errorf("CreateResource: failed to unmarshal response")
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	if len(m["parent"]) == 0 {
		return nil, status.New(codes.InvalidArgument, "invalid argument {parent does not exist}").Err()
	}
	if len(m["resource"]) > 0 {
		return nil, status.New(codes.AlreadyExists, "conflict").Err()
	}

	resource := &grbac.Resource{
		Name:   req.Resource.Name,
		Parent: req.Resource.Parent,
		Etag:   etag,
	}

	return resource, nil
}
