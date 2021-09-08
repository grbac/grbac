package graph

import (
	"context"
	"encoding/json"

	_ "embed"

	"github.com/dgraph-io/dgo/v210"
)

//go:embed data/resources.get.query.dql
var queryGetResource string

//go:embed data/resources.exists.query.dql
var queryExistsResource string

//go:embed data/resources.has_children.query.dql
var queryHasChildren string

func GetResource(ctx context.Context, txn *dgo.Txn, name string) (*Resource, error) {
	m := map[string]string{"$name": name}
	resp, err := txn.QueryWithVars(ctx, queryGetResource, m)
	if err != nil {
		return nil, err
	}

	resources := new(struct {
		Resources []*Resource `json:"resources"`
	})

	if err := json.Unmarshal(resp.Json, &resources); err != nil {
		return nil, err
	}

	if len(resources.Resources) == 0 {
		return nil, nil
	}

	return resources.Resources[0], nil
}

func ExistsResource(ctx context.Context, txn *dgo.Txn, name string) (bool, error) {
	m := map[string]string{"$name": name}
	resp, err := txn.QueryWithVars(ctx, queryExistsResource, m)
	if err != nil {
		return false, err
	}

	resources := new(struct {
		Resources []*Resource `json:"resources"`
	})

	if err := json.Unmarshal(resp.Json, &resources); err != nil {
		return false, err
	}

	return len(resources.Resources) != 0, nil
}

func HasChildren(ctx context.Context, txn *dgo.Txn, name string) (bool, error) {
	m := map[string]string{"$name": name}
	resp, err := txn.QueryWithVars(ctx, queryHasChildren, m)
	if err != nil {
		return false, err
	}

	children := new(struct {
		Resources []*Resource `json:"children"`
	})

	if err := json.Unmarshal(resp.Json, &children); err != nil {
		return false, err
	}

	return len(children.Resources) != 0, nil
}
