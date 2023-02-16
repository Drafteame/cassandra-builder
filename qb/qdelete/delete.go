package qdelete

import (
	"github.com/gocql/gocql"

	"github.com/Drafteame/cassandra-builder/qb"
	"github.com/Drafteame/cassandra-builder/qb/query"
)

type Client interface {
	Session() *gocql.Session
	Config() qb.Config
	Restart() error
	Debug() bool
	PrintFn() query.DebugPrint
}

// Query represents a Cassandra delete query. Execution should not bind any value
type Query struct {
	client Client
	table  string
	where  []query.WhereStm
	args   []interface{}
}

// New create a new delete query instance by passing a cassandra session
func New(c Client) *Query {
	return &Query{client: c}
}

// From set table where be data deleted
func (dq *Query) From(t string) *Query {
	dq.table = t
	return dq
}

// Where set where conditions that can be nested to delete data
func (dq *Query) Where(f string, op query.WhereOp, v interface{}) *Query {
	dq.where = append(dq.where, query.WhereStm{Field: f, Op: op, Value: v})
	dq.args = append(dq.args, v)
	return dq
}
