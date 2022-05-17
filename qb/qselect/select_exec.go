package qselect

import (
	"errors"
	"reflect"

	"github.com/Drafteame/cassandra-builder/qb/query"
	"github.com/gocql/gocql"
)

// One return just one result on bind action
func (q *Query) One() error {
	if q.bind == nil {
		return errors.New("nil bind is not allowed, use None() functions instead One()")
	}

	if err := query.VerifyBind(q.bind, reflect.Struct); err != nil {
		return err
	}

	sq := q.build()

	var jsonRow string

	if err := q.ctx.Session.Query(sq, q.args...).Consistency(gocql.One).Scan(&jsonRow); err != nil {
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
	if q.bind == nil {
		return errors.New("nil bind is not allowed, use None() function instead All()")
	}

	if err := query.VerifyBind(q.bind, reflect.Slice); err != nil {
		return err
	}

	sq := q.build()

	var jsonRow string

	iter := q.ctx.Session.Query(sq, q.args...).Iter()
	defer func() { _ = iter.Close() }()

	ib := reflect.Indirect(reflect.ValueOf(q.bind))

	bv := reflect.ValueOf(ib.Interface())
	bt := bv.Type().Elem()

	for iter.Scan(&jsonRow) {
		elem, err := query.BindRow([]byte(jsonRow), bt)
		if err != nil {
			return err
		}

		ib.Set(reflect.Append(ib, reflect.Indirect(elem)))
	}

	return nil
}

// None execute a qselect query without returning values
func (q *Query) None() error {
	sq := q.build()

	if err := q.ctx.Session.Query(sq, q.args...).Exec(); err != nil {
		return err
	}

	return nil
}
