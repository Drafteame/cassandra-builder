package scanner

import (
	"reflect"

	"github.com/Drafteame/cassandra-builder/qb/errors"
	"github.com/Drafteame/cassandra-builder/qb/query"
	"github.com/Drafteame/cassandra-builder/qb/runner"
	"github.com/gocql/gocql"
)

type Scanner struct {
	query     *gocql.Query
	bind      interface{}
	pageSize  int
	pageState []byte
	runner    *runner.Runner
}

func New(query *gocql.Query, pageSize int, runner *runner.Runner) *Scanner {
	return &Scanner{
		query:    query,
		pageSize: pageSize,
		runner:   runner,
	}
}

func (s Scanner) HasNextPage() bool {
	return s.pageState == nil || len(s.pageState) != 0
}

func (s *Scanner) NextPage() error {
	if s.bind == nil {
		return errors.ErrNilBinding
	}

	if err := query.VerifyBind(s.bind, reflect.Slice); err != nil {
		return err
	}

	iter := s.query.PageSize(s.pageSize).PageState(s.pageState).Iter()

	if err := s.runner.QueryWithPagination(iter, s.bind); err != nil {
		return err
	}

	// Set pastStage for next request
	s.pageState = iter.PageState()

	return nil
}

func (s *Scanner) Bind(bind interface{}) {
	s.bind = bind
}
