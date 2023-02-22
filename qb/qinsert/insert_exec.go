package qinsert

import (
	"strings"

	"github.com/Drafteame/cassandra-builder/qb/runner"
	"github.com/scylladb/gocqlx/qb"
)

// Exec execute insert query with args
func (iq *Query) Exec() error {
	run := runner.New(iq.client)
	q := iq.Build()

	if err := run.QueryNone(q, iq.args); err != nil {
		return err
	}

	return nil
}

func (iq *Query) Build() string {
	q := qb.Insert(iq.table)
	q.Columns(iq.fields...)

	queryStr, _ := q.ToCql()

	return strings.TrimSpace(queryStr)
}
