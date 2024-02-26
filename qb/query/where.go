package query

import (
	"github.com/scylladb/gocqlx/qb"
)

// WhereOp represents a where operation type
type WhereOp string

const (
	// Eq equals
	Eq WhereOp = "="

	// Ge grater or equal
	Ge WhereOp = ">="

	// Le lower or equal
	Le WhereOp = "<="

	// L lower
	L WhereOp = "<"

	// G grater
	G WhereOp = ">"

	// Tag defines a struct tag that will be used to bind row data into the given structs
	Tag = "gocql"
)

// WhereStm defines a Where condition holding filed, operation and value of the condition
type WhereStm struct {
	Field string
	Op    WhereOp
	Value any
}

// BuildWhere create a complete where statement to be used on select, delete, count and update queries
func BuildWhere(stms []WhereStm) []qb.Cmp {
	var ops []qb.Cmp

	for _, op := range stms {
		switch op.Op {
		case Eq:
			ops = append(ops, qb.Eq(op.Field))
		case Ge:
			ops = append(ops, qb.GtOrEq(op.Field))
		case Le:
			ops = append(ops, qb.LtOrEq(op.Field))
		case G:
			ops = append(ops, qb.Gt(op.Field))
		case L:
			ops = append(ops, qb.Lt(op.Field))
		}
	}

	return ops
}
