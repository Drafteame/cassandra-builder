package qb

import (
	"github.com/gocql/gocql"

	"github.com/Drafteame/cassandra-builder/qb/models"
	"github.com/Drafteame/cassandra-builder/session"
)

type client struct {
	canRestart     bool
	config         models.Config
	session        *gocql.Session
	sessionCreator session.Creator
}

var _ Client = &client{}

// NewClient creates a new cassandra client manager from config
func NewClient(conf models.Config, opts ...Option) (Client, error) {
	copts := options{sessionCreator: session.Create}

	for _, opt := range opts {
		opt(&copts)
	}

	sess, err := copts.sessionCreator(conf)
	if err != nil {
		return nil, err
	}

	return &client{
		session:        sess,
		config:         conf,
		canRestart:     true,
		sessionCreator: copts.sessionCreator,
	}, nil
}

// NewClientWithSession creates a new cassandra client manager from a given session.
func NewClientWithSession(sess *gocql.Session, conf models.Config) Client {
	return &client{
		session:    sess,
		config:     conf,
		canRestart: false,
	}
}
