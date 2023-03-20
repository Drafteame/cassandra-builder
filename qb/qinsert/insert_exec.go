package qinsert

import (
	"strings"

	"github.com/scylladb/gocqlx/qb"

	"github.com/Drafteame/cassandra-builder/qb/runner"
)

// Exec execute insert query with args
func (iq *Query) Exec() error {
	run := runner.New(iq.client)
	q := iq.build()

	if err := run.QueryNone(q, iq.args); err != nil {
		return err
	}

	return nil
}

func (iq *Query) build() string {
	q := qb.Insert(iq.table)
	q.Columns(iq.fields...)

	queryStr, _ := q.ToCql()

	return strings.TrimSpace(queryStr)
}
