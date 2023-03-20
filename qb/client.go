package qb

import (
	"github.com/gocql/gocql"

	models "github.com/Drafteame/cassandra-builder/qb/models"
	"github.com/Drafteame/cassandra-builder/qb/qcount"
	"github.com/Drafteame/cassandra-builder/qb/qdelete"
	"github.com/Drafteame/cassandra-builder/qb/qinsert"
	"github.com/Drafteame/cassandra-builder/qb/qselect"
	"github.com/Drafteame/cassandra-builder/qb/qupdate"
)

type client struct {
	canRestart bool
	config     models.Config
	session    *gocql.Session
}

var _ Client = &client{}

func (c *client) Select(f ...string) *qselect.Query {
	return qselect.New(c).Fields(f...)
}

func (c *client) Insert(f ...string) *qinsert.Query {
	return qinsert.New(c).Fields(f...)
}

func (c *client) Update(t string) *qupdate.Query {
	return qupdate.New(c).Table(t)
}

func (c *client) Delete() *qdelete.Query {
	return qdelete.New(c)
}

func (c *client) Count() *qcount.Query {
	return qcount.New(c)
}

func (c *client) Debug() bool {
	return c.config.Debug
}

func (c *client) Close() {
	c.session.Close()
}

func (c *client) Session() *gocql.Session {
	return c.session
}

func (c *client) Config() models.Config {
	return c.config
}

func (c *client) Restart() error {
	c.Close()

	session, err := createSession(c.config)
	if err != nil {
		return err
	}

	c.session = session

	return nil
}

func createSession(c models.Config) (*gocql.Session, error) {
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

	if c.Timeout != 0 {
		cluster.Timeout = c.Timeout
	}

	if c.ConnectTimeout != 0 {
		cluster.ConnectTimeout = c.ConnectTimeout
	}

	return cluster.CreateSession()
}

// NewClient creates a new cassandra client manager from config
func NewClient(conf models.Config) (Client, error) {
	session, err := createSession(conf)
	if err != nil {
		return nil, err
	}

	return &client{
		session:    session,
		config:     conf,
		canRestart: true,
	}, nil
}

// NewClientWithSession creates a new cassandra client manager from a given session.
func NewClientWithSession(session *gocql.Session, conf models.Config) Client {
	return &client{
		session:    session,
		config:     conf,
		canRestart: false,
	}
}
