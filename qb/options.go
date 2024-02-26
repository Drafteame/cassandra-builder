package qb

import "github.com/Drafteame/cassandra-builder/session"

type options struct {
	sessionCreator session.Creator
}

// Option is a function that sets some option on the client.
type Option func(*options)

// WithSessionCreator sets the session creator on the client.
func WithSessionCreator(creator session.Creator) Option {
	return func(o *options) {
		o.sessionCreator = creator
	}
}
