package qdelete

import (
	"strings"

	"github.com/scylladb/gocqlx/qb"

	"github.com/Drafteame/cassandra-builder/qb/query"
	"github.com/Drafteame/cassandra-builder/runner"
)

// Exec execute delete query and return error on failure
func (dq *Query) Exec() error {
	run := runner.New(dq.client)

	q := dq.build()

	if err := dq.client.Session().Query(q, dq.args...).Exec(); err != nil {
		return err
	}

	return nil
}

func (dq *Query) build() string {
	q := qb.Delete(dq.table)

	if len(dq.where) > 0 {
		q = q.Where(query.BuildWhere(dq.where)...)
	}

	queryStr, _ := q.ToCql()

	if dq.client.Debug() {
		dq.client.PrintFn()(queryStr, dq.args)
	}

	return strings.TrimSpace(queryStr)
}
