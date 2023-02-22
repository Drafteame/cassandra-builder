package qupdate

import (
	"strings"

	"github.com/scylladb/gocqlx/qb"

	"github.com/Drafteame/cassandra-builder/qb/query"
	"github.com/Drafteame/cassandra-builder/qb/runner"
)

// Exec run update query from builder and return an error if exists
func (uq *Query) Exec() error {
	run := runner.New(uq.client)
	q := uq.Build()

	if err := run.QueryNone(q, uq.args); err != nil {
		return err
	}

	return nil
}

func (uq *Query) Build() string {
	q := qb.Update(uq.table)

	if len(uq.fields) > 0 {
		q = q.Set(uq.fields...)
	}

	if len(uq.where) > 0 {
		if len(uq.where) > 0 {
			q = q.Where(query.BuildWhere(uq.where)...)
		}
	}

	queryStr, _ := q.ToCql()

	return strings.TrimSpace(queryStr)
}
