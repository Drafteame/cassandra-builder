package qupdate

import (
	"reflect"
	"testing"

	"github.com/Drafteame/cassandra-builder/qb/query"
)

func TestQuery_Set(t *testing.T) {
	q := &Query{}

	q.Set("col1", 1)
	q.Set("col2", 2)

	exp1 := query.Columns{"col1", "col2"}
	exp2 := []interface{}{1, 2}

	if !reflect.DeepEqual(q.fields, exp1) {
		t.Errorf("fields err: exp: %v got: %v", exp1, q.fields)
	}

	if !reflect.DeepEqual(q.args, exp2) {
		t.Errorf("args err: exp: %v got: %v", exp2, q.args)
	}
}

func TestQuery_Where(t *testing.T) {
	q := &Query{}

	q.Where("field", query.Eq, 1)
	q.Where("field2", ">=", 12)

	exp1 := []query.WhereStm{
		{
			Field: "field",
			Op:    query.Eq,
			Value: 1,
		},
		{
			Field: "field2",
			Op:    query.Ge,
			Value: 12,
		},
	}

	exp2 := []interface{}{1, 12}

	if !reflect.DeepEqual(q.where, exp1) {
		t.Errorf("fields err: exp: %v got: %v", exp1, q.where)
	}

	if !reflect.DeepEqual(q.args, exp2) {
		t.Errorf("args err: exp: %v got: %v", exp2, q.args)
	}
}
