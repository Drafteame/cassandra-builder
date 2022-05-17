package qdelete

import (
	"strings"

	"github.com/Drafteame/cassandra-builder/qb/query"
	"github.com/scylladb/gocqlx/qb"
)

// Exec execute delete query and return error on failure
func (dq *Query) Exec() error {
	q := dq.build()

	if err := dq.ctx.Session.Query(q, dq.args...).Exec(); err != nil {
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

	if dq.ctx.Debug {
		dq.ctx.PrintQuery(queryStr, dq.args)
	}

	return strings.TrimSpace(queryStr)
}
