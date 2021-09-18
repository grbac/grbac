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

func TestIntegrationAuthorize(t *testing.T) {
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
		Anonymous = "user:anonymous"

		User0 = &grbac.Subject{
			Name: "users/user-0." + uuid.New().String(),
		}
		User1 = &grbac.Subject{
			Name: "users/user-1." + uuid.New().String(),
		}
		User2 = &grbac.Subject{
			Name: "users/user-2." + uuid.New().String(),
		}
		UserNotFound = &grbac.Subject{
			Name: "users/user-?." + uuid.New().String(),
		}

		ServiceAccount0 = &grbac.Subject{
			Name: "serviceAccounts/serviceaccount-0." + uuid.New().String(),
		}
		ServiceAccount1 = &grbac.Subject{
			Name: "serviceAccounts/serviceaccount-1." + uuid.New().String(),
		}
		ServiceAccount2 = &grbac.Subject{
			Name: "serviceAccounts/serviceaccount-2." + uuid.New().String(),
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
		Group1 = &grbac.Group{
			Name: "groups/group-1." + uuid.New().String(),
			Members: []string{
				toUserMember(User1.Name),
				toServiceAccountMember(ServiceAccount1.Name),
			},
		}

		PermissionGet = &grbac.Permission{
			Name: "permissions/grbac.test.get",
		}
		PermissionCreate = &grbac.Permission{
			Name: "permissions/grbac.test.create",
		}
		PermissionDelete = &grbac.Permission{
			Name: "permissions/grbac.test.delete",
		}
		PermissionNotFound = &grbac.Permission{
			Name: "permissions/grbac.test." + uuid.New().String(),
		}

		RoleAdmin = &grbac.Role{
			Name: "roles/grbac.admin",
			Permissions: []string{
				toPermissionId(PermissionGet.Name),
				toPermissionId(PermissionCreate.Name),
				toPermissionId(PermissionDelete.Name),
			},
		}
		RoleEditor = &grbac.Role{
			Name: "roles/grbac.editor",
			Permissions: []string{
				toPermissionId(PermissionGet.Name),
				toPermissionId(PermissionCreate.Name),
			},
		}
		RoleViewer = &grbac.Role{
			Name: "roles/grbac.viewer",
			Permissions: []string{
				toPermissionId(PermissionGet.Name),
			},
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
			Parent: "@animeshon",
		}
		ResourceNotFound = &grbac.Resource{
			Name:   "//test.animeapis.com/resources/resource-?." + uuid.New().String(),
			Parent: "@animeshon",
		}

		Policy0 = &iam.Policy{
			Version: 1,
			Bindings: []*iam.Binding{
				{
					Role: RoleEditor.Name,
					Members: []string{
						toGroupMember(Group0.Name),
					},
				},
			},
		}
		Policy1 = &iam.Policy{
			Version: 1,
			Bindings: []*iam.Binding{
				{
					Role: RoleAdmin.Name,
					Members: []string{
						toGroupMember(Group0.Name),
					},
				},
				{
					Role: RoleEditor.Name,
					Members: []string{
						toUserMember(User1.Name),
						toServiceAccountMember(ServiceAccount1.Name),
					},
				},
				{
					Role: RoleViewer.Name,
					Members: []string{
						"allUsers",
					},
				},
			},
		}
		Policy2 = &iam.Policy{
			Version: 1,
			Bindings: []*iam.Binding{
				{
					Role: RoleViewer.Name,
					Members: []string{
						toGroupMember(Group0.Name),
						toGroupMember(Group1.Name),
					},
				},
			},
		}
	)

	// Create new random resources.
	_, err = server.CreateResource(context.TODO(), &grbac.CreateResourceRequest{Resource: Resource0})
	require.NoError(t, err)
	_, err = server.CreateResource(context.TODO(), &grbac.CreateResourceRequest{Resource: Resource1})
	require.NoError(t, err)
	_, err = server.CreateResource(context.TODO(), &grbac.CreateResourceRequest{Resource: Resource2})
	require.NoError(t, err)

	_, err = server.CreatePermission(context.TODO(), &grbac.CreatePermissionRequest{Permission: PermissionGet})
	require.NoError(t, err)
	_, err = server.CreatePermission(context.TODO(), &grbac.CreatePermissionRequest{Permission: PermissionCreate})
	require.NoError(t, err)
	_, err = server.CreatePermission(context.TODO(), &grbac.CreatePermissionRequest{Permission: PermissionDelete})
	require.NoError(t, err)

	_, err = server.CreateRole(context.TODO(), &grbac.CreateRoleRequest{Role: RoleAdmin})
	require.NoError(t, err)
	_, err = server.CreateRole(context.TODO(), &grbac.CreateRoleRequest{Role: RoleEditor})
	require.NoError(t, err)
	_, err = server.CreateRole(context.TODO(), &grbac.CreateRoleRequest{Role: RoleViewer})
	require.NoError(t, err)

	_, err = server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: User0})
	require.NoError(t, err)
	_, err = server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: User1})
	require.NoError(t, err)
	_, err = server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: User2})
	require.NoError(t, err)
	_, err = server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: ServiceAccount0})
	require.NoError(t, err)
	_, err = server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: ServiceAccount1})
	require.NoError(t, err)
	_, err = server.CreateSubject(context.TODO(), &grbac.CreateSubjectRequest{Subject: ServiceAccount2})
	require.NoError(t, err)

	_, err = server.CreateGroup(context.TODO(), &grbac.CreateGroupRequest{Group: Group0})
	require.NoError(t, err)
	_, err = server.CreateGroup(context.TODO(), &grbac.CreateGroupRequest{Group: Group1})
	require.NoError(t, err)

	// Set IAM polices to resources.
	_, err = server.SetIamPolicy(context.TODO(), &iam.SetIamPolicyRequest{Resource: Resource0.Name, Policy: Policy0})
	require.NoError(t, err)
	_, err = server.SetIamPolicy(context.TODO(), &iam.SetIamPolicyRequest{Resource: Resource1.Name, Policy: Policy1})
	require.NoError(t, err)
	_, err = server.SetIamPolicy(context.TODO(), &iam.SetIamPolicyRequest{Resource: Resource2.Name, Policy: Policy2})
	require.NoError(t, err)

	type T struct {
		object   string
		subject  string
		relation string
		status   codes.Code
	}

	for _, i := range []*T{
		// Test: authorization rule on non-existing resource should return permission denied.
		{ResourceNotFound.Name, User0.Name, PermissionGet.Name, codes.NotFound},
		{ResourceNotFound.Name, Anonymous, PermissionGet.Name, codes.NotFound},

		// Test: authorization rule on non-existing permission should return permission denied.
		{Resource0.Name, User0.Name, PermissionNotFound.Name, codes.PermissionDenied},
		{Resource0.Name, Anonymous, PermissionNotFound.Name, codes.PermissionDenied},

		// Test: only members of group-0 should be granted "grbac.test.create" permission on resource-0.
		{Resource0.Name, User0.Name, PermissionCreate.Name, codes.OK},
		{Resource0.Name, ServiceAccount0.Name, PermissionCreate.Name, codes.OK},

		{Resource0.Name, User1.Name, PermissionCreate.Name, codes.PermissionDenied},
		{Resource0.Name, User2.Name, PermissionCreate.Name, codes.PermissionDenied},
		{Resource0.Name, UserNotFound.Name, PermissionCreate.Name, codes.PermissionDenied},
		{Resource0.Name, ServiceAccount1.Name, PermissionCreate.Name, codes.PermissionDenied},
		{Resource0.Name, ServiceAccount2.Name, PermissionCreate.Name, codes.PermissionDenied},
		{Resource0.Name, ServiceAccountNotFound.Name, PermissionCreate.Name, codes.PermissionDenied},
		{Resource0.Name, Anonymous, PermissionCreate.Name, codes.PermissionDenied},

		// Test: only members of group-0 should be granted "grbac.test.get" permission on resource-0.
		{Resource0.Name, User0.Name, PermissionGet.Name, codes.OK},
		{Resource0.Name, ServiceAccount0.Name, PermissionGet.Name, codes.OK},

		{Resource0.Name, User1.Name, PermissionGet.Name, codes.PermissionDenied},
		{Resource0.Name, User2.Name, PermissionGet.Name, codes.PermissionDenied},
		{Resource0.Name, UserNotFound.Name, PermissionGet.Name, codes.PermissionDenied},
		{Resource0.Name, ServiceAccount1.Name, PermissionGet.Name, codes.PermissionDenied},
		{Resource0.Name, ServiceAccount2.Name, PermissionGet.Name, codes.PermissionDenied},
		{Resource0.Name, ServiceAccountNotFound.Name, PermissionGet.Name, codes.PermissionDenied},
		{Resource0.Name, Anonymous, PermissionGet.Name, codes.PermissionDenied},

		// Test: nobody should be granted "grbac.test.delete" permission on resource-0.
		{Resource0.Name, User0.Name, PermissionDelete.Name, codes.PermissionDenied},
		{Resource0.Name, User1.Name, PermissionDelete.Name, codes.PermissionDenied},
		{Resource0.Name, User2.Name, PermissionDelete.Name, codes.PermissionDenied},
		{Resource0.Name, UserNotFound.Name, PermissionDelete.Name, codes.PermissionDenied},
		{Resource0.Name, ServiceAccount0.Name, PermissionDelete.Name, codes.PermissionDenied},
		{Resource0.Name, ServiceAccount1.Name, PermissionDelete.Name, codes.PermissionDenied},
		{Resource0.Name, ServiceAccount2.Name, PermissionDelete.Name, codes.PermissionDenied},
		{Resource0.Name, ServiceAccountNotFound.Name, PermissionDelete.Name, codes.PermissionDenied},
		{Resource0.Name, Anonymous, PermissionDelete.Name, codes.PermissionDenied},

		// Test: all users should be granted "grbac.test.get" permission on resource-1.
		{Resource1.Name, User0.Name, PermissionGet.Name, codes.OK},
		{Resource1.Name, User1.Name, PermissionGet.Name, codes.OK},
		{Resource1.Name, User2.Name, PermissionGet.Name, codes.OK},
		{Resource1.Name, UserNotFound.Name, PermissionGet.Name, codes.OK},
		{Resource1.Name, ServiceAccount0.Name, PermissionGet.Name, codes.OK},
		{Resource1.Name, ServiceAccount1.Name, PermissionGet.Name, codes.OK},
		{Resource1.Name, ServiceAccount2.Name, PermissionGet.Name, codes.OK},
		{Resource1.Name, ServiceAccountNotFound.Name, PermissionGet.Name, codes.OK},
		{Resource1.Name, Anonymous, PermissionGet.Name, codes.OK},

		// Test: only members of group-0 should be granted "grbac.test.delete" permission on resource-1.
		{Resource1.Name, User0.Name, PermissionDelete.Name, codes.OK},
		{Resource1.Name, ServiceAccount0.Name, PermissionDelete.Name, codes.OK},

		{Resource1.Name, User1.Name, PermissionDelete.Name, codes.PermissionDenied},
		{Resource1.Name, User2.Name, PermissionDelete.Name, codes.PermissionDenied},
		{Resource1.Name, UserNotFound.Name, PermissionDelete.Name, codes.PermissionDenied},
		{Resource1.Name, ServiceAccount1.Name, PermissionDelete.Name, codes.PermissionDenied},
		{Resource1.Name, ServiceAccount2.Name, PermissionDelete.Name, codes.PermissionDenied},
		{Resource1.Name, ServiceAccountNotFound.Name, PermissionDelete.Name, codes.PermissionDenied},
		{Resource1.Name, Anonymous, PermissionDelete.Name, codes.PermissionDenied},

		// Test: only members of group-0 (inherited) and group-1 should be granted "grbac.test.create" permission on resource-1.
		{Resource1.Name, User0.Name, PermissionCreate.Name, codes.OK},
		{Resource1.Name, User1.Name, PermissionCreate.Name, codes.OK},
		{Resource1.Name, ServiceAccount0.Name, PermissionCreate.Name, codes.OK},
		{Resource1.Name, ServiceAccount1.Name, PermissionCreate.Name, codes.OK},

		{Resource1.Name, User2.Name, PermissionCreate.Name, codes.PermissionDenied},
		{Resource1.Name, ServiceAccount2.Name, PermissionCreate.Name, codes.PermissionDenied},
		{Resource1.Name, Anonymous, PermissionCreate.Name, codes.PermissionDenied},

		// Test: only members of group-0 and group-1 should be granted "grbac.test.get" permission on resource-2.
		{Resource2.Name, User0.Name, PermissionGet.Name, codes.OK},
		{Resource2.Name, User1.Name, PermissionGet.Name, codes.OK},
		{Resource2.Name, ServiceAccount0.Name, PermissionGet.Name, codes.OK},
		{Resource2.Name, ServiceAccount1.Name, PermissionGet.Name, codes.OK},

		{Resource2.Name, User2.Name, PermissionGet.Name, codes.PermissionDenied},
		{Resource2.Name, ServiceAccount2.Name, PermissionGet.Name, codes.PermissionDenied},
		{Resource2.Name, Anonymous, PermissionGet.Name, codes.PermissionDenied},
	} {
		subject := i.subject
		if isUser(i.subject) {
			subject = toUserMember(i.subject)
		} else if isServiceAccount(i.subject) {
			subject = toServiceAccountMember(i.subject)
		}
		_, err = server.TestIamPolicy(context.TODO(), &grbac.TestIamPolicyRequest{
			AccessTuple: &grbac.AccessTuple{
				FullResourceName: i.object,
				Principal:        subject,
				Permission:       toPermissionId(i.relation),
			},
		})

		if err == nil {
			assert.Equal(t, i.status, codes.OK)
			assert.NoError(t, err, "[%s:%s:%s]", i.object, i.relation, i.subject)
		} else {
			assert.Equal(t, i.status, status.Code(err))
			assert.Error(t, err, "[%s:%s:%s]", i.object, i.relation, i.subject)
		}
	}
}
