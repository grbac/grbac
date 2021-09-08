package graph

import (
	"context"
	"encoding/json"

	_ "embed"

	"github.com/dgraph-io/dgo/v210"
)

//go:embed data/permissions.exists.query.dql
var queryExistsPermission string

func ExistsPermission(ctx context.Context, txn *dgo.Txn, name string) (bool, error) {
	m := map[string]string{"$name": name}
	resp, err := txn.QueryWithVars(ctx, queryExistsPermission, m)
	if err != nil {
		return false, err
	}

	permissions := new(struct {
		Permissions []*Permission `json:"permissions"`
	})

	if err := json.Unmarshal(resp.Json, &permissions); err != nil {
		return false, err
	}

	return len(permissions.Permissions) != 0, nil
}
