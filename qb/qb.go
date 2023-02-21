package qb

import (
	"log"

	"github.com/gocql/gocql"

	models "github.com/Drafteame/cassandra-builder/qb/models"
	"github.com/Drafteame/cassandra-builder/qb/qcount"
	delete2 "github.com/Drafteame/cassandra-builder/qb/qdelete"
	"github.com/Drafteame/cassandra-builder/qb/qinsert"
	_select "github.com/Drafteame/cassandra-builder/qb/qselect"
	"github.com/Drafteame/cassandra-builder/qb/query"
	"github.com/Drafteame/cassandra-builder/qb/qupdate"
)

// Client is the main cassandra client abstraction to work with the database
type Client interface {
	// Select start a select query
	Select(f ...string) *_select.Query

	// Insert start a new insert query statement
	Insert(f ...string) *qinsert.Query

	// Update start an update query statement
	Update(t string) *qupdate.Query

	// Delete start a new delete query statement
	Delete() *delete2.Query

	// Count start new count query statement
	Count() *qcount.Query

	// Session return the plain session object to build some direct query
	Session() *gocql.Session

	// Debug return an assertion for debugging
	Debug() bool

	// PrintFn return the configured debug print function.
	PrintFn() query.DebugPrint

	// Restart should close and start a new connection.
	Restart() error

	// Config return current client configuration
	Config() models.Config

	// Close ends cassandra connection pool
	Close()
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
		printQuery: DefaultDebugPrint,
	}

	if conf.PrintQuery != nil {
		c.printQuery = conf.PrintQuery
	}

	return c, nil
}

// NewClientWithSession creates a new cassandra client manager from a given session.
func NewClientWithSession(session *gocql.Session, conf models.Config) (Client, error) {
	c := &client{
		session:    session,
		config:     conf,
		canRestart: false,
		printQuery: DefaultDebugPrint,
	}

	if conf.PrintQuery != nil {
		c.printQuery = conf.PrintQuery
	}

	return c, nil
}
