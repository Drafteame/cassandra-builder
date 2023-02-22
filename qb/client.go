package qb

import (
	"log"

	"github.com/gocql/gocql"

	"github.com/Drafteame/cassandra-builder/qb/qcount"
	"github.com/Drafteame/cassandra-builder/qb/qdelete"

	models "github.com/Drafteame/cassandra-builder/qb/models"
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

	session, err := getSession(c.config)
	if err != nil {
		return err
	}

	c.session = session

	return nil
}

func getSession(c models.Config) (*gocql.Session, error) {
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

// DefaultDebugPrint defines a default function that prints resultant query and arguments before being executed
// and when the Debug flag is true
func DefaultDebugPrint(q string, args []interface{}, err error) {
	if q != "" {
		log.Printf("query: %v \nargs: %v\n", q, args)
	}

	if err != nil {
		log.Println("err: ", err.Error())
	}
}

// NewClient creates a new cassandra client manager from config
func NewClient(conf models.Config) (Client, error) {
	session, err := getSession(conf)
	if err != nil {
		return nil, err
	}

	c := &client{
		session:    session,
		config:     conf,
		canRestart: true,
	}

	return c, nil
}

// NewClientWithSession creates a new cassandra client manager from a given session.
func NewClientWithSession(session *gocql.Session, conf models.Config) (Client, error) {
	c := &client{
		session:    session,
		config:     conf,
		canRestart: false,
	}

	return c, nil
}
