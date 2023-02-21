package qinsert

import (
	"reflect"
	"testing"

	"github.com/Drafteame/cassandra-builder/qb/runner/mocks"
)

func TestInsertQuery_build(t *testing.T) {
	tt := []struct {
		table   string
		fields  []string
		values  []interface{}
		res     string
		expArgs []interface{}
	}{
		{
			table:   "test1",
			fields:  []string{"col1", "col2", "col3"},
			values:  []interface{}{1, "123", true},
			expArgs: []interface{}{1, "123", true},
			res:     "INSERT INTO test1 (col1,col2,col3) VALUES (?,?,?)",
		},
		{
			table:   "test2",
			fields:  []string{"col1", "col2", "col3"},
			values:  []interface{}{"1", "123", "true"},
			expArgs: []interface{}{"1", "123", "true"},
			res:     "INSERT INTO test2 (col1,col2,col3) VALUES (?,?,?)",
		},
	}

	for _, test := range tt {
		client := mocks.NewClient(t)

		client.On("Debug").Return(false)

		q := New(client).Fields(test.fields...).Into(test.table).Values(test.values...)
		query := q.Build()

		if query != test.res {
			t.Errorf("query err: \nexp: '%v' \ngot: '%v'", test.res, query)
		}

		for i := 0; i < len(q.args); i++ {
			if !reflect.DeepEqual(q.args[i], test.expArgs[i]) {
				t.Errorf("args err: expected: (%T)%v got: (%T)%v", test.expArgs[i], test.expArgs[i], q.args[i], q.args[i])
			}
		}
	}
}
