package session

import (
	"github.com/gocql/gocql"

	"github.com/Drafteame/cassandra-builder/qb/models"
)

// Creator is a function that creates a new Cassandra session using the given configuration.
type Creator func(cfg models.Config) (*gocql.Session, error)
