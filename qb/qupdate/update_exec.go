package qupdate

import (
	"strings"

	"github.com/Drafteame/cassandra-builder/qb/query"
	"github.com/scylladb/gocqlx/qb"
)

// Exec run update query from builder and return an error if exists
func (uq *Query) Exec() error {
	q := uq.build()

	if err := uq.ctx.Session.Query(q, uq.args...).Exec(); err != nil {
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

	if uq.ctx.Debug {
		uq.ctx.PrintQuery(queryStr, uq.args)
	}

	return strings.TrimSpace(queryStr)
}
