package qinsert

import (
	"github.com/Drafteame/cassandra-builder/qb/query"
	"github.com/Drafteame/cassandra-builder/qb/runner"
)

// Query represent a Cassandra insert query. Execution should not bind any value
type Query struct {
	client runner.Client
	table  string
	fields query.Columns
	args   []interface{}
}

// New creates a new insert query by passing a cassandra session and debug options
func New(c runner.Client) *Query {
	return &Query{client: c}
}

// Fields save query fields that should be used for insert query
func (iq *Query) Fields(f ...string) *Query {
	iq.fields = f
	return iq
}

// Into set table to insert query
func (iq *Query) Into(t string) *Query {
	iq.table = t
	return iq
}

// Values set values as query arguments for insert statement
func (iq *Query) Values(v ...interface{}) *Query {
	iq.args = v
	return iq
}
