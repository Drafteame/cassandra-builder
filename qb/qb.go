package qb

import (
	"time"

	"github.com/gocql/gocql"

	"github.com/Drafteame/cassandra-builder/qb/qcount"
	delete2 "github.com/Drafteame/cassandra-builder/qb/qdelete"
	"github.com/Drafteame/cassandra-builder/qb/qinsert"
	_select "github.com/Drafteame/cassandra-builder/qb/qselect"
	"github.com/Drafteame/cassandra-builder/qb/query"
	"github.com/Drafteame/cassandra-builder/qb/qupdate"
)

type Consistency uint16

const (
	Any         Consistency = 0x00
	One         Consistency = 0x01
	Two         Consistency = 0x02
	Three       Consistency = 0x03
	Quorum      Consistency = 0x04
	All         Consistency = 0x05
	LocalQuorum Consistency = 0x06
	EachQuorum  Consistency = 0x07
	LocalOne    Consistency = 0x0A
)

// Config is the main cassandra configuration needed
type Config struct {
	Port                     int           `yaml:"port" json:"port"`
	KeyspaceName             string        `yaml:"keyspace_name" json:"keyspace_name"`
	Username                 string        `yaml:"username" json:"username"`
	Password                 string        `yaml:"password" json:"password"`
	ContactPoints            []string      `yaml:"contact_points" json:"contact_points"`
	Debug                    bool          `yaml:"debug" json:"debug"`
	ProtoVersion             int           `yaml:"proto_version" json:"proto_version"`
	Consistency              Consistency   `yaml:"consistency" json:"consistency"`
	CaPath                   string        `yaml:"ca_path" json:"ca_path"`
	DisableInitialHostLookup bool          `yaml:"disable_initial_host_lookup" json:"disable_initial_host_lookup"`
	Timeout                  time.Duration `yaml:"timeout" json:"timeout"`
	ConnectTimeout           time.Duration `yaml:"connect_timeout" json:"connect_timeout"`
	PrintQuery               query.DebugPrint
}

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

	// Close close cassandra connection pool
	Close()
}
