package runner

import (
	"reflect"

	"github.com/avast/retry-go/v4"
	"github.com/gocql/gocql"

	"github.com/Drafteame/cassandra-builder/qb/errors"
	"github.com/Drafteame/cassandra-builder/qb/models"
	"github.com/Drafteame/cassandra-builder/qb/query"
)

type Client interface {
	Session() *gocql.Session
	Config() models.Config
	Restart() error
	Debug() bool
	PrintFn() query.DebugPrint
}

type Runner struct {
	client Client
}

func (r *Runner) Query(stmt string, args []interface{}, bind interface{}) error {
	execFn := func() error {
		if r.client.Session() == nil || r.client.Session().Closed() {
			return errors.ErrClosedConnection
		}

		return r.queryAll(stmt, args, bind)
	}

	opts := r.getRetryOptions()

	if err := retry.Do(execFn, opts...); err != nil {
		return err
	}

	return nil
}

func (r *Runner) QueryOne(query string, args []interface{}) (interface{}, error) {
	var out interface{}

	execFn := func() error {
		if r.client.Session() == nil || r.client.Session().Closed() {
			return errors.ErrClosedConnection
		}

		return r.client.Session().Query(query, args...).Consistency(gocql.One).Scan(&out)
	}

	opts := r.getRetryOptions()

	if err := retry.Do(execFn, opts...); err != nil {
		return "", err
	}

	return out, nil
}

func (r *Runner) QueryNone(query string, args []interface{}) error {
	execFn := func() error {
		if r.client.Session() == nil || r.client.Session().Closed() {
			return errors.ErrClosedConnection
		}

		return r.client.Session().Query(query, args...).Exec()
	}

	opts := r.getRetryOptions()

	return retry.Do(execFn, opts...)
}

func New(c Client) *Runner {
	return &Runner{client: c}
}

func (r *Runner) getRetryOptions() []retry.Option {
	return []retry.Option{
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
}

func (r *Runner) queryAll(stmt string, args []interface{}, bind interface{}) error {
	var jsonRow string

	iter := r.client.Session().Query(stmt, args...).Iter()
	if iter == nil {
		return errors.ErrNilIterator
	}

	ib := reflect.Indirect(reflect.ValueOf(bind))

	bv := reflect.ValueOf(ib.Interface())
	bt := bv.Type().Elem()

	for iter.Scan(&jsonRow) {
		elem, err := query.BindRow([]byte(jsonRow), bt)
		if err != nil {
			return err
		}

		ib.Set(reflect.Append(ib, reflect.Indirect(elem)))
	}

	return iter.Close()
}
