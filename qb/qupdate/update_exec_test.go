package qupdate

import (
	"reflect"
	"testing"

	"github.com/Drafteame/cassandra-builder/qb/query"
)

func TestQuery_build(t *testing.T) {
	type set struct {
		field string
		value interface{}
	}

	tt := []struct {
		table   string
		set     []set
		where   []query.WhereStm
		expArgs []interface{}
		res     string
	}{
		{
			table: "test_table",
			set:   []set{},
			where: []query.WhereStm{
				{
					Field: "field1",
					Op:    query.Eq,
					Value: 0,
				},
			},
			expArgs: []interface{}{0},
			res:     "UPDATE test_table SET  WHERE field1=?",
		},
		{
			table: "test_table",
			set: []set{
				{
					field: "field1",
					value: 1,
				},
				{
					field: "field2",
					value: true,
				},
			},
			where:   []query.WhereStm{},
			expArgs: []interface{}{1, true},
			res:     "UPDATE test_table SET field1=?,field2=?",
		},
		{
			table: "test_table",
			set: []set{
				{
					field: "field1",
					value: 1,
				},
				{
					field: "field2",
					value: true,
				},
			},
			where: []query.WhereStm{
				{
					Field: "field1",
					Op:    query.Eq,
					Value: 0,
				},
			},
			expArgs: []interface{}{1, true, 0},
			res:     "UPDATE test_table SET field1=?,field2=? WHERE field1=?",
		},
	}

	for _, test := range tt {
		q := New(nil, false, nil).Table(test.table)

		for _, s := range test.set {
			q = q.Set(s.field, s.value)
		}

		for _, w := range test.where {
			q = q.Where(w.Field, w.Op, w.Value)
		}

		qs := q.build()

		if qs != test.res {
			t.Errorf("query err: \nexp: '%s' \ngot: '%s'", test.res, qs)
		}

		if !reflect.DeepEqual(test.expArgs, q.args) {
			t.Errorf("query err: \nexp: %v \ngot: %v", test.expArgs, q.args)
		}
	}
}
