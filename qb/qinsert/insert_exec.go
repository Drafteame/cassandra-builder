package qinsert

import (
	"strings"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/qb"

	"github.com/Drafteame/cassandra-builder/qb/errors"
)

// Exec execute insert query with args
func (iq *Query) Exec() error {
	q := iq.build()

	execFn := func() error {
		if r.client.Session() == nil || r.client.Session().Closed() {
			return errors.ErrClosedConnection
		}

		return r.client.Session().Query(q, iq.args...).Exec()
	}

	opts := []retry.Option{
		retry.Attempts(r.client.Config().NumRetries),
		retry.RetryIf(func(err error) bool {
			switch err {
			case gocql.ErrNoConnections, errors.ErrClosedConnection:
				return true
			default:
				return false
			}
		}),
		retry.OnRetry(func(n uint, err error) {
			errRestart := r.client.Restart()

			if r.client.Debug() {
				r.client.PrintFn()("", nil, errRestart)
			}
		}),
	}

	return retry.Do(execFn, opts...)
}

func (iq *Query) build() string {
	q := qb.Insert(iq.table)
	q.Columns(iq.fields...)

	queryStr, _ := q.ToCql()

	if iq.client.Debug() {
		iq.client.PrintFn()(queryStr, iq.args)
	}

	return strings.TrimSpace(queryStr)
}
