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

func TestIntegrationRoleCreate(t *testing.T) {
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

		Role0 = &grbac.Role{
			Name: "roles/role-0." + uuid.New().String(),
			Permissions: []string{
				toPermissionId(Permission0.Name),
			},
		}
		Role1 = &grbac.Role{
			Name: "roles/role-1." + uuid.New().String(),
			Permissions: []string{
				toPermissionId(Permission0.Name),
				toPermissionId(PermissionNotFound.Name),
			},
		}
	)

	// Create a new permission.
	_, err = server.CreatePermission(context.TODO(), &grbac.CreatePermissionRequest{Permission: Permission0})
	require.NoError(t, err)

	// Test: creation should not fail.
	role, err := server.CreateRole(context.TODO(), &grbac.CreateRoleRequest{Role: Role0})
	require.NoError(t, err)
	require.NotNil(t, role)

	assert.Equal(t, Role0.Name, role.Name)
	assert.ElementsMatch(t, Role0.Permissions, role.Permissions)
	assert.NotEmpty(t, role.Etag)

	// Test: creation with non-existing permission should fail.
	_, err = server.CreateRole(context.TODO(), &grbac.CreateRoleRequest{Role: Role1})
	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))

	// Test: creation of duplicate role should fail with already exists.
	_, err = server.CreateRole(context.TODO(), &grbac.CreateRoleRequest{Role: Role0})
	assert.Error(t, err)
	assert.Equal(t, codes.AlreadyExists, status.Code(err))

	// Test: get role should return the same role created.
	role, err = server.GetRole(context.TODO(), &grbac.GetRoleRequest{Name: Role0.Name})
	require.NoError(t, err)
	require.NotNil(t, role)

	assert.Equal(t, Role0.Name, role.Name)
	assert.ElementsMatch(t, Role0.Permissions, role.Permissions)
	assert.NotEmpty(t, role.Etag)
}

func TestIntegrationRoleDelete(t *testing.T) {
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

		Role0 = &grbac.Role{
			Name: "roles/role-0." + uuid.New().String(),
			Permissions: []string{
				toPermissionId(Permission0.Name),
			},
		}
		RoleNotFound = &grbac.Role{
			Name: "roles/role-?." + uuid.New().String(),
		}
	)

	// Create a new random role and permission.
	_, err = server.CreatePermission(context.TODO(), &grbac.CreatePermissionRequest{Permission: Permission0})
	require.NoError(t, err)

	_, err = server.CreateRole(context.TODO(), &grbac.CreateRoleRequest{Role: Role0})
	require.NoError(t, err)

	// Test: deletion of existing role should not fail.
	empty, err := server.DeleteRole(context.TODO(), &grbac.DeleteRoleRequest{Name: Role0.Name})
	assert.NoError(t, err)
	assert.NotNil(t, empty)

	// Test: get role should return 'not found' after deletion.
	_, err = server.GetRole(context.TODO(), &grbac.GetRoleRequest{Name: Role0.Name})
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))

	// Test: deletion of already deleted role should fail.
	_, err = server.DeleteRole(context.TODO(), &grbac.DeleteRoleRequest{Name: Role0.Name})
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))

	// Test: deletion of non-existing role should fail.
	_, err = server.DeleteRole(context.TODO(), &grbac.DeleteRoleRequest{Name: RoleNotFound.Name})
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestIntegrationRoleUpdate(t *testing.T) {
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
		Permission1 = &grbac.Permission{
			Name: "permissions/grbac.test." + uuid.New().String(),
		}
		PermissionNotFound = &grbac.Permission{
			Name: "permissions/grbac.test." + uuid.New().String(),
		}

		Role0 = &grbac.Role{
			Name: "roles/role-0." + uuid.New().String(),
			Permissions: []string{
				toPermissionId(Permission0.Name),
			},
		}
		RoleNotFound = &grbac.Role{
			Name: "roles/role-?." + uuid.New().String(),
		}
	)

	// Create new random roles.
	_, err = server.CreatePermission(context.TODO(), &grbac.CreatePermissionRequest{Permission: Permission0})
	require.NoError(t, err)
	_, err = server.CreatePermission(context.TODO(), &grbac.CreatePermissionRequest{Permission: Permission1})
	require.NoError(t, err)

	_, err = server.CreateRole(context.TODO(), &grbac.CreateRoleRequest{Role: Role0})
	require.NoError(t, err)

	// Test: update (replace permissions) should not fail.
	Role0.Permissions = []string{toPermissionId(Permission1.Name)}
	role, err := server.UpdateRole(context.TODO(), &grbac.UpdateRoleRequest{Role: Role0})
	require.NoError(t, err)
	require.NotNil(t, role)

	assert.Equal(t, Role0.Name, role.Name)
	assert.ElementsMatch(t, Role0.Permissions, role.Permissions)
	assert.NotEmpty(t, role.Etag)

	role, err = server.GetRole(context.TODO(), &grbac.GetRoleRequest{Name: Role0.Name})
	require.NoError(t, err)
	require.NotNil(t, role)

	assert.Equal(t, Role0.Name, role.Name)
	assert.ElementsMatch(t, Role0.Permissions, role.Permissions)
	assert.NotEmpty(t, role.Etag)

	// Test: update (add permissions) should not fail.
	Role0.Permissions = append(Role0.Permissions, toPermissionId(Permission0.Name))
	role, err = server.UpdateRole(context.TODO(), &grbac.UpdateRoleRequest{Role: Role0})
	require.NoError(t, err)
	require.NotNil(t, role)

	assert.Equal(t, Role0.Name, role.Name)
	assert.ElementsMatch(t, Role0.Permissions, role.Permissions)
	assert.NotEmpty(t, role.Etag)

	role, err = server.GetRole(context.TODO(), &grbac.GetRoleRequest{Name: Role0.Name})
	require.NoError(t, err)
	require.NotNil(t, role)

	assert.Equal(t, Role0.Name, role.Name)
	assert.ElementsMatch(t, Role0.Permissions, role.Permissions)
	assert.NotEmpty(t, role.Etag)

	// Test: update (remove all permissions) should not fail.
	Role0.Permissions = nil
	role, err = server.UpdateRole(context.TODO(), &grbac.UpdateRoleRequest{Role: Role0})
	require.NoError(t, err)
	require.NotNil(t, role)

	assert.Equal(t, Role0.Name, role.Name)
	assert.ElementsMatch(t, Role0.Permissions, role.Permissions)
	assert.NotEmpty(t, role.Etag)

	role, err = server.GetRole(context.TODO(), &grbac.GetRoleRequest{Name: Role0.Name})
	require.NoError(t, err)
	require.NotNil(t, role)

	assert.Equal(t, Role0.Name, role.Name)
	assert.ElementsMatch(t, Role0.Permissions, role.Permissions)
	assert.NotEmpty(t, role.Etag)

	// Test: update (add non-existing permission) should fail.
	Role0.Permissions = []string{toPermissionId(PermissionNotFound.Name)}
	_, err = server.UpdateRole(context.TODO(), &grbac.UpdateRoleRequest{Role: Role0})
	require.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))

	// Test: update with mutable field mask should not fail.
	Role0.Permissions = []string{toPermissionId(Permission0.Name)}
	_, err = server.UpdateRole(context.TODO(), &grbac.UpdateRoleRequest{
		Role: Role0,
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"role", "role.permissions"},
		}})
	require.NoError(t, err)

	// Test: update with immutable field mask should fail.
	_, err = server.UpdateRole(context.TODO(), &grbac.UpdateRoleRequest{
		Role: Role0,
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"role.name"},
		}})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// Test: update with invalid field mask should fail.
	_, err = server.UpdateRole(context.TODO(), &grbac.UpdateRoleRequest{
		Role: Role0,
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"<invalid>"},
		}})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// Test: update of non-existing role should fail.
	_, err = server.DeleteRole(context.TODO(), &grbac.DeleteRoleRequest{Name: RoleNotFound.Name})
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}
