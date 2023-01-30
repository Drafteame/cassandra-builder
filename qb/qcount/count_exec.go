package qcount

import (
	"strings"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/qb"

	"github.com/Drafteame/cassandra-builder/qb/query"
)

// Exec release count query an return the number of rows and a possible error
func (cq *Query) Exec() (int64, error) {
	q := cq.build()

	var count int64

	if err := cq.ctx.Session.Query(q, cq.args...).Consistency(gocql.One).Scan(&count); err != nil {
		return 0, err
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

	if cq.ctx.Debug {
		cq.ctx.PrintQuery(queryStr, cq.args)
	}

	return strings.TrimSpace(queryStr)
}
