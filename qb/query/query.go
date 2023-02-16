package query

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/gocql/gocql"

	"github.com/Drafteame/cassandra-builder/qb/errors"
)

type (
	// Columns defines a list of table column names
	Columns []string

	// Order definition for type or order on a qselect query
	Order string

	// DebugPrint defines a callback that prints query values
	DebugPrint func(q string, args []interface{}, err error)

	// Query Base definition of query
	Query struct {
		Session    *gocql.Session
		Debug      bool
		PrintQuery DebugPrint
	}
)

const (
	// Desc represents DESC order filter
	Desc Order = "DESC"

	// Asc represents ASC order filter
	Asc            Order  = "ASC"
	DatetimeLayout string = "2006-01-02 15:04:05.000Z"
)

// VerifyBind verify if an interface is bindable or not by checking it is a Ptr kind
func VerifyBind(b interface{}, k reflect.Kind) error {
	t := reflect.TypeOf(b)
	v := reflect.ValueOf(b)

	if t.Kind() != reflect.Ptr {
		return errors.ErrNoPtrBinding
	}

	str := reflect.Indirect(v).Interface()
	t = reflect.TypeOf(str)

	if t.Kind() != k {
		return errors.ErrNoStructOrSliceBinding
	}

	return nil
}

// BindMapToStruct bind values of a map into a new struct of the given reflected type Value
func BindMapToStruct(m map[string]interface{}, st reflect.Value) error {
	indVal := reflect.Indirect(st)
	indTyp := indVal.Type()

	if st.Kind() != reflect.Ptr {
		return errors.ErrNoPtrBinding
	}

	if indVal.Kind() != reflect.Struct {
		return errors.ErrNoSliceOfStructsBinding
	}

	numField := indTyp.NumField()

	for i := 0; i < numField; i++ {
		field := indTyp.Field(i)
		value := indVal.Field(i)

		tagField := field.Tag.Get(Tag)
		mv, ok := m[tagField]

		if tagField != "" && ok {
			if err := CastMapValue(mv, field.Type, value); err != nil {
				return err
			}
		}
	}

	return nil
}

// CastMapValue convert interface value into a compatible value of the reflected type
func CastMapValue(mv interface{}, t reflect.Type, v reflect.Value) error {
	newVal := reflect.Indirect(reflect.New(t)).Interface()

	switch newVal.(type) {
	case int64:
		intVal := int64(mv.(float64))
		v.SetInt(intVal)
	case float64:
		floatVal, _ := mv.(float64)
		v.SetFloat(floatVal)
	case string:
		stringVal, _ := mv.(string)
		v.SetString(stringVal)
	case bool:
		boolVal, _ := mv.(bool)
		v.SetBool(boolVal)
	case time.Time:
		t, _ := time.Parse(DatetimeLayout, mv.(string))
		v.Set(reflect.ValueOf(t))
	default:
		return fmt.Errorf("cassandra-builder: can't cast value of type %T with value %v, to type %v", mv, mv, t)
	}

	return nil
}

// BindRow bind a row of returned data on json format into a given struct type with `gocql` tags that assoc
// fields to struct members
func BindRow(jsonRow []byte, st reflect.Type) (reflect.Value, error) {
	var row map[string]interface{}
	elem := reflect.New(st)

	if err := json.Unmarshal(jsonRow, &row); err != nil {
		return reflect.Value{}, err
	}

	if err := BindMapToStruct(row, elem); err != nil {
		return reflect.Value{}, err
	}

	return elem, nil
}
