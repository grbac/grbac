package graph

import (
	"context"
	"encoding/json"

	_ "embed"

	"github.com/dgraph-io/dgo/v210"
)

//go:embed data/roles.get.query.dql
var queryGetRole string

//go:embed data/roles.exists.query.dql
var queryExistsRole string

func GetRole(ctx context.Context, txn *dgo.Txn, name string) (*Role, error) {
	m := map[string]string{"$name": name}
	resp, err := txn.QueryWithVars(ctx, queryGetRole, m)
	if err != nil {
		return nil, err
	}

	roles := new(struct {
		Roles []*Role `json:"roles"`
	})

	if err := json.Unmarshal(resp.Json, &roles); err != nil {
		return nil, err
	}

	if len(roles.Roles) == 0 {
		return nil, nil
	}

	return roles.Roles[0], nil
}

func ExistsRole(ctx context.Context, txn *dgo.Txn, name string) (bool, error) {
	m := map[string]string{"$name": name}
	resp, err := txn.QueryWithVars(ctx, queryExistsRole, m)
	if err != nil {
		return false, err
	}

	roles := new(struct {
		Roles []*Role `json:"roles"`
	})

	if err := json.Unmarshal(resp.Json, &roles); err != nil {
		return false, err
	}

	return len(roles.Roles) != 0, nil
}
