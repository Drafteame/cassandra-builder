package qupdate

import (
	"strings"

	"github.com/scylladb/gocqlx/qb"

	"github.com/Drafteame/cassandra-builder/qb/query"
)

// Exec run update query from builder and return an error if exists
func (uq *Query) Exec() error {
	q := uq.build()

	if err := uq.client.Session().Query(q, uq.args...).Exec(); err != nil {
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
