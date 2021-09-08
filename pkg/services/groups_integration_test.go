// +build integration

package services

import (
	"context"
	"os"
	"testing"

	grbac "github.com/animeapis/go-genproto/grbac/v1alpha1"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/google/uuid"
	"github.com/grbac/grbac/pkg/bootstrap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestIntegrationGroupCreate(t *testing.T) {
	endpoint := os.Getenv("INTEGRATION_TEST_DGRAPH_ENDPOINT")
	require.NotEmpty(t, endpoint)

	err := bootstrap.Schema(context.TODO(), endpoint)
	require.NoError(t, err)

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	server := &AccessControlServerImpl{
		cli:  dgo.NewDgraphClient(api.NewDgraphClient(conn)),
		conn: conn,
	}

	var (
		User0 = &grbac.Subject{
			Name: "users/user-0." + uuid.New().String(),
		}
		User1 = &grbac.Subject{
			Name: "users/user-1." + uuid.New().String(),
		}

		ServiceAccount0 = &grbac.Subject{
			Name: "serviceAccounts/serviceaccount-0." + uuid.New().String(),
		}
		ServiceAccount1 = &grbac.Subject{
			Name: "serviceAccounts/serviceaccount-1." + uuid.New().String(),
		}

		Group0 = &grbac.Group{
			Name: "groups/group-0." + uuid.New().String(),
			Members: []string{
				"allUsers",
				toUserMember(User0.Name),
				toUserMember(User1.Name),
				toServiceAccountMember(ServiceAccount0.Name),
				toServiceAccountMember(ServiceAccount1.Name),
			},
		}
		Group1 = &grbac.Group{
			Name: "groups/group-1." + uuid.New().String(),
			Members: []string{
				toGroupMember(Group0.Name),
			},
		}
		Group2 = &grbac.Group{
			Name: "groups/group-2." + uuid.New().String(),
			Members: []string{
				"allUsers",
				toUserMember(User0.Name),
				toUserMember(User1.Name),
				toServiceAccountMember(ServiceAccount0.Name),
				toServiceAccountMember(ServiceAccount1.Name),
				toGroupMember(Group0.Name),
			},
		}
		Group3 = &grbac.Group{
			Name:    "groups/group-3." + uuid.New().String(),
			Members: []string{},
		}
	)

	// Test: creation with non-existing subjects should fail.
	_, err = server.CreateGroup(context.TODO(), &grbac.CreateGroupRequest{Group: Group0})
	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))

	// Test: creation with non-existing groups should fail.
	_, err = server.CreateGroup(context.TODO(), &grbac.CreateGroupRequest{Group: Group1})
	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))

	// Test: creation with non-existing mixed members should fail.
	_, err = server.CreateGroup(context.TODO(), &grbac.CreateGroupRequest{Group: Group2})
	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))

	// Create new random subjects.
	_, err = server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: User0})
	require.NoError(t, err)

	_, err = server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: User1})
	require.NoError(t, err)

	_, err = server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: ServiceAccount0})
	require.NoError(t, err)

	_, err = server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: ServiceAccount1})
	require.NoError(t, err)

	// Test: creation (subjects only) should not fail.
	group0, err := server.CreateGroup(context.TODO(), &grbac.CreateGroupRequest{Group: Group0})
	require.NoError(t, err)
	require.NotNil(t, group0)

	assert.Equal(t, Group0.Name, group0.Name)
	assert.ElementsMatch(t, Group0.Members, group0.Members)
	assert.NotEmpty(t, group0.Etag)

	// Test: creation (groups only) should not fail.
	group1, err := server.CreateGroup(context.TODO(), &grbac.CreateGroupRequest{Group: Group1})
	require.NoError(t, err)
	require.NotNil(t, group1)

	assert.Equal(t, Group1.Name, group1.Name)
	assert.ElementsMatch(t, Group1.Members, group1.Members)
	assert.NotEmpty(t, group1.Etag)

	// Test: creation (mixed members) should not fail.
	group2, err := server.CreateGroup(context.TODO(), &grbac.CreateGroupRequest{Group: Group2})
	require.NoError(t, err)
	require.NotNil(t, group2)

	assert.Equal(t, Group2.Name, group2.Name)
	assert.ElementsMatch(t, Group2.Members, group2.Members)
	assert.NotEmpty(t, group2.Etag)

	// Test: creation (no members) should not fail.
	group3, err := server.CreateGroup(context.TODO(), &grbac.CreateGroupRequest{Group: Group3})
	require.NoError(t, err)
	require.NotNil(t, group3)

	assert.Equal(t, Group3.Name, group3.Name)
	assert.Empty(t, group3.Members)
	assert.NotEmpty(t, group3.Etag)

	// Test: creation of duplicate group should fail with already exists.
	_, err = server.CreateGroup(context.TODO(), &grbac.CreateGroupRequest{Group: Group0})
	assert.Error(t, err)
	assert.Equal(t, codes.AlreadyExists, status.Code(err))

	// Test: get group (mixed members) should return the same group created.
	group, err := server.GetGroup(context.TODO(), &grbac.GetGroupRequest{Name: Group2.Name})
	require.NoError(t, err)
	require.NotNil(t, group)

	assert.Equal(t, Group2.Name, group.Name)
	assert.ElementsMatch(t, Group2.Members, group.Members)
	assert.NotEmpty(t, group.Etag)

	// Test: get group (no members) should return the same group created.
	group, err = server.GetGroup(context.TODO(), &grbac.GetGroupRequest{Name: Group3.Name})
	require.NoError(t, err)
	require.NotNil(t, group)

	assert.Equal(t, Group3.Name, group.Name)
	assert.Empty(t, group.Members)
	assert.NotEmpty(t, group.Etag)
}

func TestIntegrationGroupDelete(t *testing.T) {
	endpoint := os.Getenv("INTEGRATION_TEST_DGRAPH_ENDPOINT")
	require.NotEmpty(t, endpoint)

	err := bootstrap.Schema(context.TODO(), endpoint)
	require.NoError(t, err)

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	server := &AccessControlServerImpl{
		cli:  dgo.NewDgraphClient(api.NewDgraphClient(conn)),
		conn: conn,
	}

	var (
		User0 = &grbac.Subject{
			Name: "users/user-0." + uuid.New().String(),
		}

		ServiceAccount0 = &grbac.Subject{
			Name: "serviceAccounts/serviceaccount-0." + uuid.New().String(),
		}

		Group0 = &grbac.Group{
			Name: "groups/group-0." + uuid.New().String(),
			Members: []string{
				"allUsers",
				toUserMember(User0.Name),
				toServiceAccountMember(ServiceAccount0.Name),
			},
		}
		GroupNotFound = &grbac.Group{
			Name: "groups/group-?." + uuid.New().String(),
		}
	)

	// Create new random group and subjects.
	_, err = server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: User0})
	require.NoError(t, err)

	_, err = server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: ServiceAccount0})
	require.NoError(t, err)

	_, err = server.CreateGroup(context.TODO(), &grbac.CreateGroupRequest{Group: Group0})
	require.NoError(t, err)

	// Test: deletion of existing resource with no children should not fail.
	empty, err := server.DeleteGroup(context.TODO(), &grbac.DeleteGroupRequest{Name: Group0.Name})
	assert.NoError(t, err)
	assert.NotNil(t, empty)

	// Test: get resource should return 'not found' after deletion.
	_, err = server.GetGroup(context.TODO(), &grbac.GetGroupRequest{Name: Group0.Name})
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))

	// Test: deletion of already deleted resource should fail.
	_, err = server.DeleteGroup(context.TODO(), &grbac.DeleteGroupRequest{Name: Group0.Name})
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))

	// Test: deletion of non-existing resource should fail.
	_, err = server.DeleteGroup(context.TODO(), &grbac.DeleteGroupRequest{Name: GroupNotFound.Name})
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestIntegrationGroupUpdate(t *testing.T) {
	endpoint := os.Getenv("INTEGRATION_TEST_DGRAPH_ENDPOINT")
	require.NotEmpty(t, endpoint)

	err := bootstrap.Schema(context.TODO(), endpoint)
	require.NoError(t, err)

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	server := &AccessControlServerImpl{
		cli:  dgo.NewDgraphClient(api.NewDgraphClient(conn)),
		conn: conn,
	}

	var (
		User0 = &grbac.Subject{
			Name: "users/user-0." + uuid.New().String(),
		}
		UserNotFound = &grbac.Subject{
			Name: "users/user-?." + uuid.New().String(),
		}

		ServiceAccount0 = &grbac.Subject{
			Name: "serviceAccounts/serviceaccount-0." + uuid.New().String(),
		}
		ServiceAccountNotFound = &grbac.Subject{
			Name: "serviceAccounts/serviceaccount-?." + uuid.New().String(),
		}

		Group0 = &grbac.Group{
			Name: "groups/group-0." + uuid.New().String(),
			Members: []string{
				toUserMember(User0.Name),
				toServiceAccountMember(ServiceAccount0.Name),
			},
		}
		GroupNotFound = &grbac.Group{
			Name: "groups/group-?." + uuid.New().String(),
		}
	)

	// Create new random group and subjects.
	_, err = server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: User0})
	require.NoError(t, err)

	_, err = server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: ServiceAccount0})
	require.NoError(t, err)

	_, err = server.CreateGroup(context.TODO(), &grbac.CreateGroupRequest{Group: Group0})
	require.NoError(t, err)

	// Test: update (add existing subjects) should not fail.
	Group0.Members = append(Group0.Members,
		"allUsers",
	)

	group, err := server.UpdateGroup(context.TODO(), &grbac.UpdateGroupRequest{Group: Group0})
	require.NoError(t, err)
	require.NotNil(t, group)

	assert.Equal(t, Group0.Name, group.Name)
	assert.ElementsMatch(t, Group0.Members, group.Members)
	assert.NotEmpty(t, group.Etag)

	group, err = server.GetGroup(context.TODO(), &grbac.GetGroupRequest{Name: Group0.Name})
	require.NoError(t, err)
	require.NotNil(t, group)

	assert.Equal(t, Group0.Name, group.Name)
	assert.ElementsMatch(t, Group0.Members, group.Members)
	assert.NotEmpty(t, group.Etag)

	// Test: update (add non-existing subjects) should fail.
	Group0.Members = append(Group0.Members,
		toUserMember(UserNotFound.Name),
		toServiceAccountMember(ServiceAccountNotFound.Name),
	)

	_, err = server.UpdateGroup(context.TODO(), &grbac.UpdateGroupRequest{Group: Group0})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// Test: update (remove subjects) should not fail.
	Group0.Members = nil
	group, err = server.UpdateGroup(context.TODO(), &grbac.UpdateGroupRequest{Group: Group0})
	require.NoError(t, err)
	require.NotNil(t, group)

	assert.Equal(t, Group0.Name, group.Name)
	assert.ElementsMatch(t, Group0.Members, group.Members)
	assert.NotEmpty(t, group.Etag)

	group, err = server.GetGroup(context.TODO(), &grbac.GetGroupRequest{Name: Group0.Name})
	require.NoError(t, err)
	require.NotNil(t, group)

	assert.Equal(t, Group0.Name, group.Name)
	assert.ElementsMatch(t, Group0.Members, group.Members)
	assert.NotEmpty(t, group.Etag)

	// Test: update with mutable field mask should not fail.
	_, err = server.UpdateGroup(context.TODO(), &grbac.UpdateGroupRequest{
		Group: Group0,
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"group", "group.members"},
		}})
	require.NoError(t, err)

	// Test: update with immutable field mask should fail.
	_, err = server.UpdateGroup(context.TODO(), &grbac.UpdateGroupRequest{
		Group: Group0,
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"group.name"},
		}})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// Test: update with invalid field mask should fail.
	_, err = server.UpdateGroup(context.TODO(), &grbac.UpdateGroupRequest{
		Group: Group0,
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"<invalid>"},
		}})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// Test: update of a self-referencing group should fail.
	Group0.Members = []string{Group0.Name}
	_, err = server.UpdateGroup(context.TODO(), &grbac.UpdateGroupRequest{Group: Group0})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// Test: update of non-existing resource should fail.
	_, err = server.DeleteGroup(context.TODO(), &grbac.DeleteGroupRequest{Name: GroupNotFound.Name})
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}
