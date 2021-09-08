package graph

import (
	"context"
	"encoding/json"

	_ "embed"

	"github.com/dgraph-io/dgo/v210"
)

//go:embed data/subjects.exists.query.dql
var queryExistsSubject string

func ExistsSubject(ctx context.Context, txn *dgo.Txn, name string) (bool, error) {
	m := map[string]string{"$name": name}
	resp, err := txn.QueryWithVars(ctx, queryExistsSubject, m)
	if err != nil {
		return false, err
	}

	subjects := new(struct {
		Subjects []*Subject `json:"subjects"`
	})

	if err := json.Unmarshal(resp.Json, &subjects); err != nil {
		return false, err
	}

	return len(subjects.Subjects) != 0, nil
}
