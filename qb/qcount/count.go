package qcount

import (
	"github.com/Drafteame/cassandra-builder/qb/query"
	"github.com/gocql/gocql"
)

// Query create new select count query
type Query struct {
	ctx            query.Query
	table          string
	column         string
	where          []query.WhereStm
	allowFiltering bool
	args           []interface{}
}

// New create a new count query instance by passing a cassandra session
func New(s *gocql.Session, d bool, dp query.DebugPrint) *Query {
	return &Query{ctx: query.Query{
		Session:    s,
		Debug:      d,
		PrintQuery: dp,
	}}
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
