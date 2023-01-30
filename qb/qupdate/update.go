package qupdate

import (
	"github.com/gocql/gocql"

	"github.com/Drafteame/cassandra-builder/qb/query"
)

// Query represent a Cassandra update query. Execution should not bind any value
type Query struct {
	ctx    query.Query
	table  string
	fields query.Columns
	args   []interface{}
	where  []query.WhereStm
}

// New create a new update query by passing a cassandra session and the affected table
func New(s *gocql.Session, d bool, dp query.DebugPrint) *Query {
	return &Query{ctx: query.Query{
		Session:    s,
		Debug:      d,
		PrintQuery: dp,
	}}
}

// Table set the table name to affect with the update query
func (uq *Query) Table(t string) *Query {
	uq.table = t
	return uq
}

// Set save field and corresponding value to update
func (uq *Query) Set(f string, v interface{}) *Query {
	uq.fields = append(uq.fields, f)
	uq.args = append(uq.args, v)
	return uq
}

// Where adds single where conditional. If more are needed, concatenate more calls to this functions
func (uq *Query) Where(f string, op query.WhereOp, v interface{}) *Query {
	uq.where = append(uq.where, query.WhereStm{Field: f, Op: op, Value: v})
	uq.args = append(uq.args, v)
	return uq
}
