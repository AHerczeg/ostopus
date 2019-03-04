package query

import (
	"github.com/sirupsen/logrus"
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

func NewQueryStore() localQueryStore {
	var qs localQueryStore
	qs.queries = make(map[string]string)

	for _, q := range queries {
		qs.AddQueries(q)
	}

	return qs
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
