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
)

func TestIntegrationPermissionCreate(t *testing.T) {
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
		Permission0 = &grbac.Permission{
			Name: "permissions/grbac.test." + uuid.New().String(),
		}
		PermissionInvalid = &grbac.Permission{
			Name: "permissions/" + uuid.New().String(),
		}
	)

	// Test: creation should not fail.
	user0, err := server.CreatePermission(context.TODO(), &grbac.CreatePermissionRequest{Permission: Permission0})
	require.NoError(t, err)
	require.NotNil(t, user0)

	assert.Equal(t, Permission0.Name, user0.Name)

	// Test: creation with invalid format should fail.
	_, err = server.CreatePermission(context.TODO(), &grbac.CreatePermissionRequest{Permission: PermissionInvalid})
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// Test: creation of duplicate permission should fail with already exists.
	_, err = server.CreatePermission(context.TODO(), &grbac.CreatePermissionRequest{Permission: Permission0})
	assert.Error(t, err)
	assert.Equal(t, codes.AlreadyExists, status.Code(err))
}

func TestIntegrationPermissionDelete(t *testing.T) {
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
		Permission0 = &grbac.Permission{
			Name: "permissions/grbac.test." + uuid.New().String(),
		}
		PermissionNotFound = &grbac.Permission{
			Name: "permissions/grbac.test." + uuid.New().String(),
		}
	)

	// Create a new random permission.
	_, err = server.CreatePermission(context.TODO(), &grbac.CreatePermissionRequest{Permission: Permission0})
	require.NoError(t, err)

	// Test: deletion of existing permission should not fail.
	empty, err := server.DeletePermission(context.TODO(), &grbac.DeletePermissionRequest{Name: Permission0.Name})
	require.NoError(t, err)
	assert.NotNil(t, empty)

	// Test: deletion of deleted permission should fail.
	_, err = server.DeletePermission(context.TODO(), &grbac.DeletePermissionRequest{Name: Permission0.Name})
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))

	// Test: deletion of non-existing permission should fail.
	_, err = server.DeletePermission(context.TODO(), &grbac.DeletePermissionRequest{Name: PermissionNotFound.Name})
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}
