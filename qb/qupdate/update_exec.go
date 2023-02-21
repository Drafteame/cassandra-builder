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
	q := uq.build()

	if err := run.QueryNone(q, uq.args); err != nil {
		if uq.client.Debug() {
			uq.client.PrintFn()(q, uq.args, err)
		}

		return err
	}

	return nil
}

func (uq *Query) build() string {
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

	if uq.client.Debug() {
		uq.client.PrintFn()(queryStr, uq.args, nil)
	}

	return strings.TrimSpace(queryStr)
}
