package query

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

var (
	queries = []map[string]string{
		standardQueries,
	}
)

type QueryStore interface {
	GetQuery(name string) (string, bool)
	HasQuery(name string) bool
	AddQueries(queries map[string]string)
	AddQuery(name, query string)
}

type localQueryStore struct {
	queries map[string]string
}

func NewLocalQueryStore() *localQueryStore {
	var qs localQueryStore
	qs.queries = make(map[string]string)

	for _, q := range queries {
		qs.AddQueries(q)
	}

	return &qs
}

func (qs *localQueryStore) GetQuery(name string) (string, bool) {
	query, ok := qs.queries[name]
	return query, ok
}

func (qs *localQueryStore) HasQuery(name string) bool {
	_, ok := qs.queries[name]
	return ok
}

func (qs *localQueryStore) AddQueries(queries map[string]string) {
	for k, v := range queries {
		qs.AddQuery(k, v)
	}
}

func (qs *localQueryStore) AddQuery(name, query string) {
	if _, ok := qs.queries[name]; ok {
		logrus.Info("A query with name already exists in query store", "name", name)
	} else {
		qs.queries[name] = query
	}
}

func Foo(i int) int {
	return i
}


type MockQueryStore struct {
	mock.Mock
}

func (m *MockQueryStore) GetQuery(name string) (string, bool) {
	args := m.Called(name)
	return args.String(0), args.Bool(1)
}

func (m *MockQueryStore) HasQuery(name string) bool {
	args := m.Called(name)
	return args.Bool(0)
}

func (m *MockQueryStore) AddQueries(queries map[string]string) {
	m.Called(queries)
}

func (m *MockQueryStore) AddQuery(name, query string) {
	m.Called(name, query)
}
