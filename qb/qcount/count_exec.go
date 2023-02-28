package qcount

import (
	"strings"

	"github.com/scylladb/gocqlx/qb"

	"github.com/Drafteame/cassandra-builder/qb/query"
	"github.com/Drafteame/cassandra-builder/qb/runner"
)

// Exec release count query an return the number of rows and a possible error
func (cq *Query) Exec() (int64, error) {
	run := runner.New(cq.client)

	q := cq.build()

	return run.QueryCount(q, cq.args)
}

func (cq *Query) build() string {
	q := qb.Select(cq.table).Count(cq.column)

	if len(cq.where) > 0 {
		q = q.Where(query.BuildWhere(cq.where)...)
	}

	if cq.allowFiltering {
		q = q.AllowFiltering()
	}

	queryStr, _ := q.ToCql()

	return strings.TrimSpace(queryStr)
}
