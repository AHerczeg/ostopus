package query

import (
	"OStopus/tentacle/os"
	"fmt"
	"github.com/sirupsen/logrus"
)

type QueryHandler struct {
	queryStore localQueryStore
	osHandler  os.OSHandler
}

func NewQueryHandler(store localQueryStore, os os.OSHandler) QueryHandler {
	return QueryHandler{queryStore: store, osHandler: os}
}

func (qh QueryHandler) RunSavedQuery(name string) (ResultDTO, error) {
	fmt.Printf("Running query \"%s\"", name)
	query, err := qh.fetchQuery(name)
	if err != nil {
		return ResultDTO{}, err
	}
	return qh.executeQuery(query)
}

func (qh QueryHandler) RunCustomQuery(query string) (ResultDTO, error) {
	logrus.Info("Running query", query)
	return qh.executeQuery(query)
}

func (qh QueryHandler) fetchQuery(name string) (string, error) {
	if query, ok := qh.queryStore.GetQuery(name); !ok {
		return "", fmt.Errorf("missing query")
	} else {
		return query, nil
	}
}

func (qh QueryHandler) executeQuery(query string) (ResultDTO, error) {
	response, err := qh.osHandler.Execute(query)
	if err != nil {
		return ResultDTO{}, err
	}
	var result ResultDTO
	result.UnmarshalArguments(response.Bytes())
	return result, nil
}
