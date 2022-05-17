package qcount

import (
	"reflect"
	"testing"

	"github.com/gocql/gocql"

	"github.com/Drafteame/cassandra-builder/qb/query"
)

func TestNew(t *testing.T) {
	s := &gocql.Session{}
	q := New(s, false, nil)

	if !reflect.DeepEqual(q.ctx.Session, s) {
		t.Errorf("associated session is different")
	}
}

func TestQuery_From(t *testing.T) {
	q := &Query{}

	q.From("test")
	if q.table != "test" {
		t.Errorf("exp: test got: %v", q.table)
	}

	q.From("test2")
	if q.table != "test2" {
		t.Errorf("exp: test2 got: %v", q.table)
	}
}

func TestQuery_Column(t *testing.T) {
	q := &Query{}

	q.Column("field")
	if q.column != "field" {
		t.Errorf("exp: field got: %v", q.column)
	}

	q.Column("field2")
	if q.column != "field2" {
		t.Errorf("exp: field2 got: %v", q.column)
	}
}

func TestQuery_Where(t *testing.T) {
	tt := []struct {
		field   string
		op      query.WhereOp
		value   interface{}
		expArgs []interface{}
		expStm  []query.WhereStm
	}{
		{
			field:   "field1",
			op:      query.Eq,
			value:   nil,
			expArgs: []interface{}{nil},
			expStm: []query.WhereStm{
				{
					Field: "field1",
					Op:    query.Eq,
					Value: nil,
				},
			},
		},
		{
			field:   "field2",
			op:      query.G,
			value:   1,
			expArgs: []interface{}{1},
			expStm: []query.WhereStm{
				{
					Field: "field2",
					Op:    query.G,
					Value: 1,
				},
			},
		},
		{
			field:   "field3",
			op:      query.L,
			value:   1,
			expArgs: []interface{}{1},
			expStm: []query.WhereStm{
				{
					Field: "field3",
					Op:    query.L,
					Value: 1,
				},
			},
		},
	}

	for _, test := range tt {
		q := &Query{}

		q.Where(test.field, test.op, test.value)

		if !reflect.DeepEqual(q.args, test.expArgs) {
			t.Errorf("exp: %v got: %v", test.expArgs, q.args)
		}

		if !reflect.DeepEqual(q.where, test.expStm) {
			t.Errorf("exp: %v got: %v", test.expStm, q.where)
		}
	}
}
