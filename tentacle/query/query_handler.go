package query

import (
	"OStopus/tentacle/os"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"regexp"
)

var (
	jsonRegex = regexp.MustCompile("{\\s*(\"[^\"]*\":\"[^\"]*\"\\s*)+(,\\s*(\"[^\"]*\":\"[^\"]*\")\\s*)*}")
)

type QueryHandler struct {
	queryStore QueryStore
	osHandler  os.OSHandler
}

func NewQueryHandler(store QueryStore, os os.OSHandler) *QueryHandler {
	return &QueryHandler{queryStore: store, osHandler: os}
}

func (qh QueryHandler) RunSavedQuery(name string) (ResultDTO, error) {
	logrus.Info("Running query: ", name)
	query, err := qh.fetchQuery(name)
	if err != nil {
		return ResultDTO{}, err
	}
	return qh.executeQuery(query)
}

func (qh QueryHandler) RunCustomQuery(query string) (ResultDTO, error) {
	logrus.Info("Running query: ", query)
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
	result.UnmarshalArguments(cleanJSON(response))
	return result, nil
}

func cleanJSON(out bytes.Buffer) []byte {
	jsonRegex.Longest()
	cleanedJSON := jsonRegex.Find(out.Bytes())
	fmt.Println(json.Valid(cleanedJSON))
	fmt.Println(cleanedJSON)
	if !json.Valid(cleanedJSON) {
		return []byte{}
	}
	return cleanedJSON
}
