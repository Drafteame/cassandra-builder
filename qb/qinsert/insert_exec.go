package qinsert

import (
	"strings"

	"github.com/scylladb/gocqlx/qb"
)

// Exec execute insert query with args
func (iq *Query) Exec() error {
	q := iq.build()

	if err := iq.ctx.Session.Query(q, iq.args...).Exec(); err != nil {
		return err
	}

	return nil
}

func (iq *Query) build() string {
	q := qb.Insert(iq.table)
	q.Columns(iq.fields...)

	queryStr, _ := q.ToCql()

	if iq.ctx.Debug {
		iq.ctx.PrintQuery(queryStr, iq.args)
	}

	return strings.TrimSpace(queryStr)
}
