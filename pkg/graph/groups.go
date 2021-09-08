package graph

import (
	"context"
	"encoding/json"

	_ "embed"

	"github.com/dgraph-io/dgo/v210"
)

//go:embed data/groups.get.query.dql
var queryGetGroup string

//go:embed data/groups.exists.query.dql
var queryExistsGroup string

func GetGroup(ctx context.Context, txn *dgo.Txn, name string) (*Group, error) {
	m := map[string]string{"$name": name}
	resp, err := txn.QueryWithVars(ctx, queryGetGroup, m)
	if err != nil {
		return nil, err
	}

	groups := new(struct {
		Groups []*Group `json:"groups"`
	})

	if err := json.Unmarshal(resp.Json, &groups); err != nil {
		return nil, err
	}

	if len(groups.Groups) == 0 {
		return nil, nil
	}

	return groups.Groups[0], nil
}

func ExistsGroup(ctx context.Context, txn *dgo.Txn, name string) (bool, error) {
	m := map[string]string{"$name": name}
	resp, err := txn.QueryWithVars(ctx, queryExistsGroup, m)
	if err != nil {
		return false, err
	}

	groups := new(struct {
		Groups []*Group `json:"groups"`
	})

	if err := json.Unmarshal(resp.Json, &groups); err != nil {
		return false, err
	}

	return len(groups.Groups) != 0, nil
}
