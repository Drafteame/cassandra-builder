package qinsert

import (
	"strings"

	"github.com/Drafteame/cassandra-builder/qb/runner"
	"github.com/scylladb/gocqlx/qb"
)

// Exec execute insert query with args
func (iq *Query) Exec() error {
	run := runner.New(iq.client)
	q := iq.build()

	if err := run.QueryNone(q, iq.args); err != nil {
		if iq.client.Debug() {
			iq.client.PrintFn()(q, iq.args, err)
		}

		return err
	}

	return nil
}

func (iq *Query) build() string {
	q := qb.Insert(iq.table)
	q.Columns(iq.fields...)

	queryStr, _ := q.ToCql()

	if iq.client.Debug() {
		iq.client.PrintFn()(queryStr, iq.args, nil)
	}

	return strings.TrimSpace(queryStr)
}
