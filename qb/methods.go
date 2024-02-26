package qb

import (
	"github.com/gocql/gocql"

	"github.com/Drafteame/cassandra-builder/qb/models"
	"github.com/Drafteame/cassandra-builder/qb/qcount"
	"github.com/Drafteame/cassandra-builder/qb/qdelete"
	"github.com/Drafteame/cassandra-builder/qb/qinsert"
	"github.com/Drafteame/cassandra-builder/qb/qselect"
	"github.com/Drafteame/cassandra-builder/qb/qupdate"
)

// Select creates a new select query.
func (c *client) Select(f ...string) *qselect.Query {
	return qselect.New(c).Fields(f...)
}

// Insert creates a new insert query.
func (c *client) Insert(f ...string) *qinsert.Query {
	return qinsert.New(c).Fields(f...)
}

// Update creates a new update query.
func (c *client) Update(t string) *qupdate.Query {
	return qupdate.New(c).Table(t)
}

// Delete creates a new delete query.
func (c *client) Delete() *qdelete.Query {
	return qdelete.New(c)
}

// Count creates a new count query.
func (c *client) Count() *qcount.Query {
	return qcount.New(c)
}

// Debug returns the debug flag.
func (c *client) Debug() bool {
	return c.config.Debug
}

// Close closes the session.
func (c *client) Close() {
	if c.session == nil {
		return
	}

	c.session.Close()
}

// Session returns the session from gocql driver
func (c *client) Session() *gocql.Session {
	return c.session
}

// Config returns the configuration used to create the session.
func (c *client) Config() models.Config {
	return c.config
}

// Restart restarts the session.
func (c *client) Restart() error {
	c.Close()

	sess, err := c.sessionCreator(c.config)
	if err != nil {
		return err
	}

	c.session = sess

	return nil
}
