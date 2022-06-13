package qb

import (
	"github.com/Drafteame/cassandra-builder/qb/qcount"
	"github.com/Drafteame/cassandra-builder/qb/qdelete"
	"github.com/Drafteame/cassandra-builder/qb/qinsert"
	"github.com/Drafteame/cassandra-builder/qb/qselect"
	"github.com/Drafteame/cassandra-builder/qb/query"
	"github.com/Drafteame/cassandra-builder/qb/qupdate"
	"github.com/gocql/gocql"
)

type client struct {
	session    *gocql.Session
	debug      bool
	printQuery query.DebugPrint
}

// NewClient creates a new cassandra client manager from config
func NewClient(conf Config) (Client, error) {
	session, err := getSession(conf)
	if err != nil {
		return nil, err
	}

	c := &client{session: session, debug: conf.Debug, printQuery: query.DefaultDebugPrint}

	if conf.PrintQuery != nil {
		c.printQuery = conf.PrintQuery
	}

	return c, nil
}

var _ Client = &client{}

func (c *client) Select(f ...string) *qselect.Query {
	return qselect.New(c.session, c.debug, c.printQuery).Fields(f...)
}

func (c *client) Insert(f ...string) *qinsert.Query {
	return qinsert.New(c.session, c.debug, c.printQuery).Fields(f...)
}

func (c *client) Update(t string) *qupdate.Query {
	return qupdate.New(c.session, c.debug, c.printQuery).Table(t)
}

func (c *client) Delete() *qdelete.Query {
	return qdelete.New(c.session, c.debug, c.printQuery)
}

func (c *client) Count() *qcount.Query {
	return qcount.New(c.session, c.debug, c.printQuery)
}

// Close finish cassandra session
func (c *client) Close() {
	c.session.Close()
}

func (c *client) Session() *gocql.Session {
	return c.session
}

func getSession(c Config) (*gocql.Session, error) {
	cluster := gocql.NewCluster(c.ContactPoints...)
	cluster.Keyspace = c.KeyspaceName
	cluster.Consistency = gocql.Consistency(c.Consistency)
	cluster.ProtoVersion = c.ProtoVersion

	if c.Port != 0 {
		cluster.Port = c.Port
	}

	if c.DisableInitialHostLookup {
		cluster.DisableInitialHostLookup = c.DisableInitialHostLookup
	}

	if c.Username != "" && c.Password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: c.Username,
			Password: c.Password,
		}
	}

	if c.CaPath != "" {
		cluster.SslOpts = &gocql.SslOptions{
			CaPath: c.CaPath,
		}
	}

	return cluster.CreateSession()
}
