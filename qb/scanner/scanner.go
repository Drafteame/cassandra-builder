package scanner

import (
	"github.com/Drafteame/cassandra-builder/qb/runner"
)

type Scanner struct {
	stmt      string
	args      []interface{}
	bind      interface{}
	pageSize  int
	pageState []byte
	runner    *runner.Runner
}

func New(stmt string, args []interface{}, bind interface{}, pageSize int, runner *runner.Runner) *Scanner {
	return &Scanner{
		stmt:      stmt,
		args:      args,
		bind:      bind,
		pageSize:  pageSize,
		pageState: nil,
		runner:    runner,
	}
}

func (s Scanner) HasNextPage() bool {
	return s.pageState != nil && len(s.pageState) != 0
}

func (s *Scanner) NextPage() error {
	iter := s.runner.NewQuery(s.stmt, s.args).PageSize(s.pageSize).PageState(s.pageState).Iter()

	if err := s.runner.QueryWithPagination(iter, s.bind); err != nil {
		return err
	}

	// Set pastStage for next request
	s.pageState = iter.PageState()

	return nil
}
