package qcount

import (
	"github.com/Drafteame/cassandra-builder/qb/query"
	"github.com/Drafteame/cassandra-builder/qb/runner"
)

// Query create new select count query
type Query struct {
	client         runner.Client
	table          string
	column         string
	where          []query.WhereStm
	allowFiltering bool
	args           []interface{}
}

// New create a new count query instance by passing a cassandra session
func New(c runner.Client) *Query {
	return &Query{client: c}
}

// Column set count column of the query
func (cq *Query) Column(c string) *Query {
	cq.column = c
	return cq
}

// From set table for count query
func (cq *Query) From(t string) *Query {
	cq.table = t
	return cq
}

// AllowFiltering sets a ALLOW FILTERING clause on the query.
func (cq *Query) AllowFiltering() *Query {
	cq.allowFiltering = true
	return cq
}

// Where adds single where conditional. If more are needed, concatenate more calls to this functions
func (cq *Query) Where(f string, op query.WhereOp, v interface{}) *Query {
	cq.where = append(cq.where, query.WhereStm{Field: f, Op: op, Value: v})
	cq.args = append(cq.args, v)
	return cq
}
