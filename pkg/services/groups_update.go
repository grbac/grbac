package services

import (
	"context"
	"encoding/base64"
	"errors"
	"text/template"

	_ "embed"

	grbac "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"github.com/dgraph-io/dgo/v210"
	"github.com/grbac/grbac/pkg/fieldmask"
	"github.com/grbac/grbac/pkg/graph"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:embed data/groups/groups.update.query.go.tmpl
var queryUpdateGroup string

//go:embed data/groups/groups.update.set.go.tmpl
var setUpdateGroup string

//go:embed data/groups/groups.update.delete.go.tmpl
var deleteUpdateGroup string

var templateQueryUpdateGroup = template.Must(
	template.New("QueryUpdateGroup").Funcs(defaultFuncMap).Parse(queryUpdateGroup),
)

var templateSetUpdateGroup = template.Must(
	template.New("SetUpdateGroup").Funcs(defaultFuncMap).Parse(setUpdateGroup),
)

var templateDeleteUpdateGroup = template.Must(
	template.New("DeleteUpdateGroup").Funcs(defaultFuncMap).Parse(deleteUpdateGroup),
)

func (s *AccessControlServerImpl) validateUpdateGroup(ctx context.Context, txn *dgo.Txn, req *grbac.UpdateGroupRequest) error {
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

	// The update field mask must contain valid paths.
	for _, path := range req.GetUpdateMask().GetPaths() {
		switch path {
		case "group", "group.members":
		default:
			return status.New(codes.InvalidArgument, "invalid argument {invalid field mask}").Err()
		}
	}

	// The members must all exist and must have a valid type.
	for _, m := range req.Group.Members {
		memberFound, err := false, error(nil)
		if isGroupMember(m) {
			if toGroupName(m) == req.Group.Name {
				return status.New(codes.InvalidArgument, "invalid argument {self-containing groups are forbidden}").Err()
			}

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
			logrus.WithError(err).Errorf("UpdateGroup: failed to query group members")
			return status.New(codes.Internal, "internal error").Err()
		}

		if !memberFound {
			return status.New(codes.InvalidArgument, "invalid argument {member does not exist}").Err()
		}
	}

	// The group must exist.
	groupFound, err := graph.ExistsGroup(ctx, txn, req.Group.Name)
	if err != nil {
		logrus.WithError(err).Errorf("UpdateGroup: failed to query group")
		return status.New(codes.Internal, "internal error").Err()
	}

	if !groupFound {
		return status.New(codes.NotFound, "not found").Err()
	}

	return nil
}

// UpdateGroup updates a group with a field mask.
func (s *AccessControlServerImpl) UpdateGroup(ctx context.Context, req *grbac.UpdateGroupRequest) (*grbac.Group, error) {
	txn := s.cli.NewTxn()
	if err := s.validateUpdateGroup(ctx, txn, req); err != nil {
		return nil, err
	}

	// TODO: etag should be generated according to the data structure.
	etag := []byte("TODO")

	fieldmask := fieldmask.NewFieldMask(req.GetUpdateMask())

	data := struct {
		Group     *grbac.Group
		FieldMask func(string) bool
		ETag      string
	}{
		Group:     req.GetGroup(),
		FieldMask: fieldmask.Contains,
		ETag:      base64.StdEncoding.EncodeToString(etag),
	}

	if err := s.update(ctx, txn, templateQueryUpdateGroup, templateSetUpdateGroup, templateDeleteUpdateGroup, data); err != nil {
		if errors.Is(err, dgo.ErrAborted) {
			return nil, status.New(codes.Aborted, "transaction has been aborted").Err()
		}

		logrus.WithError(err).Errorf("UpdateGroup: failed to execute dgraph call")
		return nil, status.New(codes.Internal, "internal error").Err()
	}

	// TODO: merge missing fields (not included in the update mask) with the group in dgraph.
	group := &grbac.Group{
		Name:    req.Group.Name,
		Members: req.Group.Members,
		Etag:    etag,
	}

	return group, nil
}
