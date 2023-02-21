package qcount

import (
	"strings"

	"github.com/scylladb/gocqlx/qb"

	"github.com/Drafteame/cassandra-builder/qb/errors"
	"github.com/Drafteame/cassandra-builder/qb/query"
	"github.com/Drafteame/cassandra-builder/qb/runner"
)

// Exec release count query an return the number of rows and a possible error
func (cq *Query) Exec() (int64, error) {
	run := runner.New(cq.client)

	q := cq.build()

	out, err := run.QueryOne(q, cq.args)
	if err != nil {
		return 0, err
	}

	count, ok := out.(int64)
	if !ok {
		return 0, errors.ErrParsing
	}

	return count, nil
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

	if cq.client.Debug() {
		cq.client.PrintFn()(queryStr, cq.args, nil)
	}

	return strings.TrimSpace(queryStr)
}
