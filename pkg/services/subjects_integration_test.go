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

func TestIntegrationSubjectCreate(t *testing.T) {
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
	)

	// Test: creation (user) should not fail.
	user0, err := server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: User0})
	require.NoError(t, err)
	require.NotNil(t, user0)

	assert.Equal(t, User0.Name, user0.Name)

	// Test: creation (serviceAccount) should not fail.
	serviceAccount, err := server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: ServiceAccount0})
	require.NoError(t, err)
	require.NotNil(t, serviceAccount)

	assert.Equal(t, ServiceAccount0.Name, serviceAccount.Name)

	// Test: creation of duplicate subject should fail with already exists.
	_, err = server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: User0})
	assert.Error(t, err)
	assert.Equal(t, codes.AlreadyExists, status.Code(err))

	_, err = server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: ServiceAccount0})
	assert.Error(t, err)
	assert.Equal(t, codes.AlreadyExists, status.Code(err))
}

func TestIntegrationSubjectDelete(t *testing.T) {
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
		Subject0 = &grbac.Subject{
			Name: "users/user-0." + uuid.New().String(),
		}
		SubjectNotFound = &grbac.Subject{
			Name: "users/user-?." + uuid.New().String(),
		}
	)

	// Create a new random subject.
	_, err = server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: Subject0})
	require.NoError(t, err)

	// Test: deletion of existing subject should not fail.
	empty, err := server.DeleteSubject(context.TODO(), &grbac.DeleteSubjectRequest{Name: Subject0.Name})
	require.NoError(t, err)
	assert.NotNil(t, empty)

	// Test: deletion of deleted subject should fail.
	_, err = server.DeleteSubject(context.TODO(), &grbac.DeleteSubjectRequest{Name: Subject0.Name})
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))

	// Test: deletion of non-existing subject should fail.
	_, err = server.DeleteSubject(context.TODO(), &grbac.DeleteSubjectRequest{Name: SubjectNotFound.Name})
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}
