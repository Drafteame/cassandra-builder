package qinsert

import (
	"reflect"
	"testing"
)

func TestInsertQuery_Into(t *testing.T) {
	q := &Query{}

	q.Into("test_table")

	if q.table != "test_table" {
		t.Errorf("expected: test_table got: %s", q.table)
	}
}

func TestInsertQuery_Values(t *testing.T) {
	q := &Query{}

	values := []interface{}{1, "asd", true, 12.34}

	q.Values(values...)

	for i := 0; i < len(values); i++ {
		if !reflect.DeepEqual(values[i], q.args[i]) {
			t.Errorf("expected: (%T)%v got: (%T)%v", values[i], values[i], q.args[i], q.args[i])
		}
	}
}
