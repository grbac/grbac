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

func TestIntegrationResourceCreate(t *testing.T) {
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
		ResourceNotFound = &grbac.Resource{
			Name:   "//test.animeapis.com/resources/resource-?." + uuid.New().String(),
			Parent: "@animeshon",
		}
		Resource0 = &grbac.Resource{
			Name:   "//test.animeapis.com/resources/resource-0." + uuid.New().String(),
			Parent: "@animeshon",
		}
		Resource1 = &grbac.Resource{
			Name:   "//test.animeapis.com/resources/resource-1." + uuid.New().String(),
			Parent: Resource0.Name,
		}
		Resource2 = &grbac.Resource{
			Name:   "//test.animeapis.com/resources/resource-2." + uuid.New().String(),
			Parent: ResourceNotFound.Name,
		}
		Resource3 = &grbac.Resource{
			Name: "//test.animeapis.com/resources/resource-3." + uuid.New().String(),
		}
	)

	// Test: creation should not fail.
	resource0, err := server.CreateResource(context.TODO(), &grbac.CreateResourceRequest{Resource: Resource0})
	require.NoError(t, err)
	require.NotNil(t, resource0)

	assert.Equal(t, Resource0.Name, resource0.Name)
	assert.Equal(t, Resource0.Parent, resource0.Parent)
	assert.NotEmpty(t, resource0.Etag)

	// Test: creation with existing parent should not fail.
	resource1, err := server.CreateResource(context.TODO(), &grbac.CreateResourceRequest{Resource: Resource1})
	require.NoError(t, err)
	require.NotNil(t, resource1)

	assert.Equal(t, Resource1.Name, resource1.Name)
	assert.Equal(t, Resource1.Parent, resource1.Parent)
	assert.NotEmpty(t, resource1.Etag)

	// Test: creation with non-existing parent should fail.
	_, err = server.CreateResource(context.TODO(), &grbac.CreateResourceRequest{Resource: Resource2})
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	// Test: creation without parent should fail.
	_, err = server.CreateResource(context.TODO(), &grbac.CreateResourceRequest{Resource: Resource3})
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	// Test: creation of duplicate resource should fail with already exists.
	_, err = server.CreateResource(context.TODO(), &grbac.CreateResourceRequest{Resource: Resource0})
	assert.Error(t, err)
	assert.Equal(t, codes.AlreadyExists, status.Code(err))

	// Test: get resource should return the same resource created.
	resource, err := server.GetResource(context.TODO(), &grbac.GetResourceRequest{Name: Resource1.Name})
	require.NoError(t, err)
	require.NotNil(t, resource)

	assert.Equal(t, Resource1.Name, resource.Name)
	assert.Equal(t, Resource1.Parent, resource.Parent)
	assert.NotEmpty(t, resource.Etag)
}

func TestIntegrationResourceDelete(t *testing.T) {
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
		Resource0 = &grbac.Resource{
			Name:   "//test.animeapis.com/resources/resource-0." + uuid.New().String(),
			Parent: "@animeshon",
		}
		Resource1 = &grbac.Resource{
			Name:   "//test.animeapis.com/resources/resource-1." + uuid.New().String(),
			Parent: Resource0.Name,
		}
		ResourceNotFound = &grbac.Resource{
			Name:   "//test.animeapis.com/resources/resource-?." + uuid.New().String(),
			Parent: "@animeshon",
		}
	)

	// Create new random resources.
	_, err = server.CreateResource(context.TODO(), &grbac.CreateResourceRequest{Resource: Resource0})
	require.NoError(t, err)

	_, err = server.CreateResource(context.TODO(), &grbac.CreateResourceRequest{Resource: Resource1})
	require.NoError(t, err)

	// Test: deletion of existing resource with children should fail.
	_, err = server.DeleteResource(context.TODO(), &grbac.DeleteResourceRequest{Name: Resource0.Name})
	assert.Error(t, err)
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))

	// Test: deletion of existing resource with no children should not fail.
	empty, err := server.DeleteResource(context.TODO(), &grbac.DeleteResourceRequest{Name: Resource1.Name})
	assert.NoError(t, err)
	assert.NotNil(t, empty)

	empty, err = server.DeleteResource(context.TODO(), &grbac.DeleteResourceRequest{Name: Resource0.Name})
	assert.NoError(t, err)
	assert.NotNil(t, empty)

	// Test: get resource should return 'not found' after deletion.
	_, err = server.GetResource(context.TODO(), &grbac.GetResourceRequest{Name: Resource0.Name})
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))

	_, err = server.GetResource(context.TODO(), &grbac.GetResourceRequest{Name: Resource1.Name})
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))

	// Test: deletion of already deleted resource should fail.
	_, err = server.DeleteResource(context.TODO(), &grbac.DeleteResourceRequest{Name: Resource0.Name})
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))

	// Test: deletion of non-existing resource should fail.
	_, err = server.DeleteResource(context.TODO(), &grbac.DeleteResourceRequest{Name: ResourceNotFound.Name})
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}
