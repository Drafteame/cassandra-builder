package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/gocql/gocql"
)

type (
	// Columns defines a list of table column names
	Columns []string

	// Order definition for type or order on a qselect query
	Order string

	// DebugPrint defines a callback that prints query values
	DebugPrint func(q string, args []interface{})

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
	Asc Order = "ASC"
)

// DefaultDebugPrint defines a default function that prints resultant query and arguments before being executed
// and when the Debug flag is true
func DefaultDebugPrint(q string, args []interface{}) {
	log.Printf("query: %v \nargs: %v\n", q, args)
}

// VerifyBind verify if an interface is bindable or not by checking it is a Ptr kind
func VerifyBind(b interface{}, k reflect.Kind) error {
	t := reflect.TypeOf(b)
	v := reflect.ValueOf(b)

	if t.Kind() != reflect.Ptr {
		return errors.New("bind value should be a pointer")
	}

	str := reflect.Indirect(v).Interface()
	t = reflect.TypeOf(str)

	if t.Kind() != k {
		return errors.New("bind value needs to be a struct to get one value, or a slice of structs to get many")
	}

	return nil
}

// BindMapToStruct bind values of a map into a new struct of the given reflected type Value
func BindMapToStruct(m map[string]interface{}, st reflect.Value) error {
	indVal := reflect.Indirect(st)
	indTyp := indVal.Type()

	if st.Kind() != reflect.Ptr {
		return errors.New("bind should be a slice of struct pointers with `gocql` tags")
	}

	if indVal.Kind() != reflect.Struct {
		return errors.New("bind should be a slice of struct pointers with `gocql` tags - 2")
	}

	numField := indTyp.NumField()

	for i := 0; i < numField; i++ {
		field := indTyp.Field(i)
		value := indVal.Field(i)

		tagField := field.Tag.Get(Tag)
		mv, ok := m[tagField]

		if tagField != "" && ok {
			fmt.Printf("map: %T field: %v\n", mv, field.Type)

			err := CastMapValue(mv, field.Type, value)
			if err != nil {
				fmt.Printf("err: %v\n", err)
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
	default:
		return fmt.Errorf("can't cast value of type %T with value %v, to type %v", mv, mv, t)
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
