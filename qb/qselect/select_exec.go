package qselect

import (
	"reflect"

	"github.com/Drafteame/cassandra-builder/qb/errors"
	"github.com/Drafteame/cassandra-builder/qb/query"
	"github.com/Drafteame/cassandra-builder/qb/runner"
	"github.com/Drafteame/cassandra-builder/qb/scanner"
)

// One return just one result on bind action
func (q *Query) One() error {
	if q.bind == nil {
		return errors.ErrNilBinding
	}

	if err := query.VerifyBind(q.bind, reflect.Struct); err != nil {
		return err
	}

	sq := q.build()

	run := runner.New(q.client)

	jsonRow, err := run.QueryOne(sq, q.args)
	if err != nil {
		return err
	}

	ib := reflect.Indirect(reflect.ValueOf(q.bind))

	bv := reflect.ValueOf(ib.Interface())
	bt := bv.Type()

	elem, err := query.BindRow([]byte(jsonRow), bt)
	if err != nil {
		return err
	}

	ib.Set(reflect.Indirect(elem))

	return nil
}

// All return all rows on bind action. Bind should be a slice of structs
func (q *Query) All() error {
	run := runner.New(q.client)
	if q.bind == nil {
		return errors.ErrNilBinding
	}

	if err := query.VerifyBind(q.bind, reflect.Slice); err != nil {
		return err
	}

	sq := q.build()

	return run.Query(sq, q.args, q.bind)
}

func (q *Query) Paginated(pageSize int) (scanner.Scanner, error) {
	run := runner.New(q.client)
	if q.bind == nil {
		return scanner.Scanner{}, errors.ErrNilBinding
	}

	if err := query.VerifyBind(q.bind, reflect.Slice); err != nil {
		return scanner.Scanner{}, err
	}

	sq := q.build()

	scanner := scanner.New(sq, q.args, q.bind, pageSize, run)

	return *scanner, nil
}
