package qcount

import (
	"reflect"
	"testing"

	"github.com/gocql/gocql"

	"github.com/Drafteame/cassandra-builder/qb/query"
)

func TestQuery_build(t *testing.T) {
	tt := []struct {
		table          string
		column         string
		where          []query.WhereStm
		expQuery       string
		expArgs        []interface{}
		allowFiltering bool
	}{
		{
			table:  "test",
			column: "*",
			where: []query.WhereStm{
				{
					Field: "field1",
					Op:    query.Eq,
					Value: true,
				},
				{
					Field: "field2",
					Op:    query.G,
					Value: 123,
				},
			},
			expArgs:  []interface{}{true, 123},
			expQuery: "SELECT count(*) FROM test WHERE field1=? AND field2>?",
		},
		{
			table:  "test",
			column: "id",
			where: []query.WhereStm{
				{
					Field: "field1",
					Op:    query.Eq,
					Value: true,
				},
				{
					Field: "field2",
					Op:    query.Ge,
					Value: 123,
				},
			},
			expArgs:  []interface{}{true, 123},
			expQuery: "SELECT count(id) FROM test WHERE field1=? AND field2>=?",
		},
		{
			table:          "test",
			column:         "id",
			allowFiltering: true,
			where: []query.WhereStm{
				{
					Field: "field1",
					Op:    query.Eq,
					Value: true,
				},
				{
					Field: "field2",
					Op:    query.Ge,
					Value: 123,
				},
			},
			expArgs:  []interface{}{true, 123},
			expQuery: "SELECT count(id) FROM test WHERE field1=? AND field2>=? ALLOW FILTERING",
		},
	}

	for i, test := range tt {
		q := New(&gocql.Session{}, false, nil).From(test.table).Column(test.column)

		for _, w := range test.where {
			q = q.Where(w.Field, w.Op, w.Value)
		}

		if test.allowFiltering {
			q = q.AllowFiltering()
		}

		queryStr := q.build()

		if test.expQuery != queryStr {
			t.Errorf("test case: %v\nexp: '%v' \ngot: '%v'", i, test.expQuery, queryStr)
		}

		if !reflect.DeepEqual(q.args, test.expArgs) {
			t.Errorf("test case: %v\nexp: %v got: %v", i, test.expArgs, q.args)
		}
	}
}
