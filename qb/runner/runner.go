package runner

import (
	"reflect"

	"github.com/avast/retry-go/v4"
	"github.com/gocql/gocql"

	"github.com/Drafteame/cassandra-builder/qb/errors"
	"github.com/Drafteame/cassandra-builder/qb/models"
	"github.com/Drafteame/cassandra-builder/qb/query"
)

//go:generate mockery --name=Client --filename=client.go --structname=Client --output=../test/mocks --outpkg=mocks

type Client interface {
	Session() *gocql.Session
	Config() models.Config
	Restart() error
	Debug() bool
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

func (r *Runner) QueryWithPagination(iter *gocql.Iter, bind interface{}) error {
	execFn := func() error {
		if r.client.Session() == nil || r.client.Session().Closed() {
			return errors.ErrClosedConnection
		}

		return r.queryWithPagination(iter, bind)
	}

	opts := r.getRetryOptions()

	if err := retry.Do(execFn, opts...); err != nil {
		return err
	}

	return nil
}

func (r *Runner) QueryCount(query string, args []interface{}) (int64, error) {
	var count int64

	execFn := func() error {
		if r.client.Session() == nil || r.client.Session().Closed() {
			return errors.ErrClosedConnection
		}

		consistency := r.client.Config().Consistency

		if err := r.client.Session().Query(query, args...).Consistency(gocql.Consistency(consistency)).Scan(&count); err != nil {
			if err == gocql.ErrNoConnections {
				return err
			}

			return retry.Unrecoverable(err)
		}

		return nil
	}

	opts := r.getRetryOptions()

	if err := retry.Do(execFn, opts...); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Runner) QueryOne(query string, args []interface{}) (string, error) {
	var jsonRow string

	execFn := func() error {
		if r.client.Session() == nil || r.client.Session().Closed() {
			return errors.ErrClosedConnection
		}

		consistency := r.client.Config().Consistency

		if err := r.client.Session().Query(query, args...).Consistency(gocql.Consistency(consistency)).Scan(&jsonRow); err != nil {
			if err == gocql.ErrNoConnections {
				return err
			}

			return retry.Unrecoverable(err)
		}

		return nil
	}

	opts := r.getRetryOptions()

	if err := retry.Do(execFn, opts...); err != nil {
		return "", err
	}

	return jsonRow, nil
}

func (r *Runner) QueryNone(query string, args []interface{}) error {
	execFn := func() error {
		if r.client.Session() == nil || r.client.Session().Closed() {
			return errors.ErrClosedConnection
		}

		if err := r.client.Session().Query(query, args...).Exec(); err != nil {
			if err == gocql.ErrNoConnections {
				return err
			}

			return retry.Unrecoverable(err)
		}

		return nil
	}

	opts := r.getRetryOptions()

	return retry.Do(execFn, opts...)
}

func New(c Client) *Runner {
	return &Runner{client: c}
}

func (r Runner) NewQuery(stmt string, args []interface{}) *gocql.Query {
	return r.client.Session().Query(stmt, args...)
}

func (r *Runner) getRetryOptions() []retry.Option {
	return []retry.Option{
		retry.Attempts(r.client.Config().NumRetries),
		retry.OnRetry(func(n uint, err error) {

			_ = r.client.Restart()

			//TODO: handle error
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

	if err := iter.Close(); err != nil {
		if err == gocql.ErrNoConnections {
			return err
		}

		return retry.Unrecoverable(err)
	}

	return nil
}

func (r *Runner) queryWithPagination(iter *gocql.Iter, bind interface{}) error {
	var jsonRow string

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

	if err := iter.Close(); err != nil {
		if err == gocql.ErrNoConnections {
			return err
		}

		return retry.Unrecoverable(err)
	}

	return nil
}
