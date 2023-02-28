package errors

import "errors"

var (
	ErrNilBinding              = errors.New("cassandra-builder: nil bind is not allowed")
	ErrNoPtrBinding            = errors.New("cassandra-builder: bind value should be a pointer")
	ErrNoStructOrSliceBinding  = errors.New("cassandra-builder: bind value should be a struct or slice")
	ErrNoSliceOfStructsBinding = errors.New("cassandra-builder: bind value should be a slice of structs")
	ErrClosedConnection        = errors.New("cassandra-builder: can execute on closed connection")
	ErrNilIterator             = errors.New("cassandra-builder: nil iterator is not allowed")
	ErrParsing                 = errors.New("cassandra-builder: error parsing row")
)
