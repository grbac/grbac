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
	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestIntegrationSetIamPolicy(t *testing.T) {
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
			Name: "users/" + uuid.New().String(),
		}
		User1 = &grbac.Subject{
			Name: "users/" + uuid.New().String(),
		}

		ServiceAccount0 = &grbac.Subject{
			Name: "serviceAccounts/" + uuid.New().String(),
		}
		ServiceAccount1 = &grbac.Subject{
			Name: "serviceAccounts/" + uuid.New().String(),
		}

		Group0 = &grbac.Group{
			Name: "groups/" + uuid.New().String(),
			Members: []string{
				toUserMember(User0.Name),
				toServiceAccountMember(ServiceAccount0.Name),
			},
		}
		Group1 = &grbac.Group{
			Name: "groups/" + uuid.New().String(),
		}

		Permission0 = &grbac.Permission{
			Name: "permissions/grbac.test." + uuid.New().String(),
		}
		Permission1 = &grbac.Permission{
			Name: "permissions/grbac.test." + uuid.New().String(),
		}

		Role0 = &grbac.Role{
			Name: "roles/" + uuid.New().String(),
			Permissions: []string{
				toPermissionId(Permission0.Name),
			},
		}
		Role1 = &grbac.Role{
			Name: "roles/" + uuid.New().String(),
			Permissions: []string{
				toPermissionId(Permission0.Name),
				toPermissionId(Permission1.Name),
			},
		}
		Role2 = &grbac.Role{
			Name: "roles/" + uuid.New().String(),
		}

		Resource0 = &grbac.Resource{
			Name:   "//test.animeapis.com/resources/" + uuid.New().String(),
			Parent: "@animeshon",
		}
		Resource1 = &grbac.Resource{
			Name:   "//test.animeapis.com/resources/" + uuid.New().String(),
			Parent: "@animeshon",
		}

		Policy0 = &iam.Policy{
			Version: 1,
			Bindings: []*iam.Binding{
				{
					Role: Role0.Name,
					Members: []string{
						toUserMember(User0.Name),
						toServiceAccountMember(ServiceAccount0.Name),
						toGroupMember(Group0.Name),
					},
				},
				{
					Role: Role1.Name,
					Members: []string{
						toUserMember(User0.Name),
					},
				},
			},
		}
	)

	// Create new random resources.
	_, err = server.CreateResource(context.TODO(), &grbac.CreateResourceRequest{Resource: Resource0})
	require.NoError(t, err)

	_, err = server.CreatePermission(context.TODO(), &grbac.CreatePermissionRequest{Permission: Permission0})
	require.NoError(t, err)

	_, err = server.CreatePermission(context.TODO(), &grbac.CreatePermissionRequest{Permission: Permission1})
	require.NoError(t, err)

	_, err = server.CreateRole(context.TODO(), &grbac.CreateRoleRequest{Role: Role0})
	require.NoError(t, err)

	_, err = server.CreateRole(context.TODO(), &grbac.CreateRoleRequest{Role: Role1})
	require.NoError(t, err)

	_, err = server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: User0})
	require.NoError(t, err)

	_, err = server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: ServiceAccount0})
	require.NoError(t, err)

	_, err = server.CreateGroup(context.TODO(), &grbac.CreateGroupRequest{Group: Group0})
	require.NoError(t, err)

	// Test: newly created resource should have an empty policy.
	policy, err := server.GetIamPolicy(context.TODO(), &iam.GetIamPolicyRequest{Resource: Resource0.Name})
	require.NoError(t, err)
	require.NotNil(t, policy)
	require.Empty(t, policy.Bindings)
	require.Empty(t, policy.Etag)
	require.Empty(t, policy.Version)

	// Test: get policy should return 'not found' if the resource doesn't exist.
	_, err = server.GetIamPolicy(context.TODO(), &iam.GetIamPolicyRequest{Resource: Resource1.Name})
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))

	// Test: setting a valid resource policy should not fail.
	policy, err = server.SetIamPolicy(context.TODO(), &iam.SetIamPolicyRequest{
		Resource: Resource0.Name,
		Policy:   Policy0,
	})
	require.NoError(t, err)
	require.NotNil(t, policy)
	require.Equal(t, Policy0.Bindings, policy.Bindings)
	require.Equal(t, Policy0.Version, policy.Version)
	require.NotEmpty(t, policy.Etag)

	// Test: get resource should return the same resource created.
	policy, err = server.GetIamPolicy(context.TODO(), &iam.GetIamPolicyRequest{Resource: Resource0.Name})
	require.NoError(t, err)
	require.NotNil(t, policy)
	require.Equal(t, Policy0.Bindings, policy.Bindings)
	require.Equal(t, Policy0.Version, policy.Version)
	require.NotEmpty(t, policy.Etag)

	// Test: setting an invalid (no policy) resource policy should fail.
	_, err = server.SetIamPolicy(context.TODO(), &iam.SetIamPolicyRequest{
		Resource: Resource0.Name,
	})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// Test: setting an invalid (no resource name) resource policy should fail.
	_, err = server.SetIamPolicy(context.TODO(), &iam.SetIamPolicyRequest{
		Policy: &iam.Policy{
			Version: 1,
			Bindings: []*iam.Binding{
				{
					Role: Role0.Name,
					Members: []string{
						toUserMember(User0.Name),
						toServiceAccountMember(ServiceAccount0.Name),
						toGroupMember(Group0.Name),
					},
				},
			},
		},
	})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// Test: setting an invalid (non-existing resource) resource policy should fail.
	_, err = server.SetIamPolicy(context.TODO(), &iam.SetIamPolicyRequest{
		Resource: Resource1.Name,
		Policy: &iam.Policy{
			Version: 1,
			Bindings: []*iam.Binding{
				{
					Role: Role0.Name,
					Members: []string{
						toUserMember(User0.Name),
						toServiceAccountMember(ServiceAccount0.Name),
						toGroupMember(Group0.Name),
					},
				},
			},
		},
	})
	require.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))

	// Test: setting an invalid (unsupported version) resource policy should fail.
	_, err = server.SetIamPolicy(context.TODO(), &iam.SetIamPolicyRequest{
		Resource: Resource0.Name,
		Policy: &iam.Policy{
			Version: 5,
			Bindings: []*iam.Binding{
				{
					Role: Role0.Name,
					Members: []string{
						toUserMember(User0.Name),
						toServiceAccountMember(ServiceAccount0.Name),
						toGroupMember(Group0.Name),
					},
				},
			},
		},
	})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// Test: setting an invalid (no role) resource policy should fail.
	_, err = server.SetIamPolicy(context.TODO(), &iam.SetIamPolicyRequest{
		Resource: Resource0.Name,
		Policy: &iam.Policy{
			Version: 1,
			Bindings: []*iam.Binding{
				{
					Members: []string{
						toUserMember(User0.Name),
						toServiceAccountMember(ServiceAccount0.Name),
						toGroupMember(Group0.Name),
					},
				},
			},
		},
	})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// Test: setting an invalid (non-existing role) resource policy should fail.
	_, err = server.SetIamPolicy(context.TODO(), &iam.SetIamPolicyRequest{
		Resource: Resource0.Name,
		Policy: &iam.Policy{
			Version: 1,
			Bindings: []*iam.Binding{
				{
					Role: Role2.Name,
					Members: []string{
						toUserMember(User0.Name),
						toServiceAccountMember(ServiceAccount0.Name),
						toGroupMember(Group0.Name),
					},
				},
			},
		},
	})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// Test: setting an invalid (non-existing user) resource policy should fail.
	_, err = server.SetIamPolicy(context.TODO(), &iam.SetIamPolicyRequest{
		Resource: Resource0.Name,
		Policy: &iam.Policy{
			Version: 1,
			Bindings: []*iam.Binding{
				{
					Role: Role0.Name,
					Members: []string{
						toUserMember(User1.Name),
						toServiceAccountMember(ServiceAccount0.Name),
						toGroupMember(Group0.Name),
					},
				},
			},
		},
	})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// Test: setting an invalid (non-existing service account) resource policy should fail.
	_, err = server.SetIamPolicy(context.TODO(), &iam.SetIamPolicyRequest{
		Resource: Resource0.Name,
		Policy: &iam.Policy{
			Version: 1,
			Bindings: []*iam.Binding{
				{
					Role: Role0.Name,
					Members: []string{
						toUserMember(User0.Name),
						toServiceAccountMember(ServiceAccount1.Name),
						toGroupMember(Group0.Name),
					},
				},
			},
		},
	})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// Test: setting an invalid (non-existing group) resource policy should fail.
	_, err = server.SetIamPolicy(context.TODO(), &iam.SetIamPolicyRequest{
		Resource: Resource0.Name,
		Policy: &iam.Policy{
			Version: 1,
			Bindings: []*iam.Binding{
				{
					Role: Role0.Name,
					Members: []string{
						toUserMember(User0.Name),
						toServiceAccountMember(ServiceAccount1.Name),
						toGroupMember(Group1.Name),
					},
				},
			},
		},
	})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// Test: setting an invalid (no members) resource policy should fail.
	_, err = server.SetIamPolicy(context.TODO(), &iam.SetIamPolicyRequest{
		Resource: Resource0.Name,
		Policy: &iam.Policy{
			Version: 1,
			Bindings: []*iam.Binding{
				{
					Role: Role0.Name,
				},
			},
		},
	})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// Test: setting an invalid (repeated role) resource policy should fail.
	_, err = server.SetIamPolicy(context.TODO(), &iam.SetIamPolicyRequest{
		Resource: Resource0.Name,
		Policy: &iam.Policy{
			Version: 1,
			Bindings: []*iam.Binding{
				{
					Role: Role0.Name,
					Members: []string{
						toUserMember(User0.Name),
					},
				}, {
					Role: Role0.Name,
					Members: []string{
						toUserMember(ServiceAccount0.Name),
					},
				},
			},
		},
	})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}
