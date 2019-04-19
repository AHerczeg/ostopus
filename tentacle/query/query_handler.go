package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"ostopus/tentacle/os"
	"regexp"
)

var (
	queryHandler Handler
	jsonRegex = regexp.MustCompile("{\\s*(\"[^\"]*\":\"[^\"]*\"\\s*)+(,\\s*(\"[^\"]*\":\"[^\"]*\")\\s*)*}")
)

type Handler struct {
	queryStore QueryStore
	osHandler  os.Handler
}

func InitQueryHandler(store QueryStore, os os.Handler) {
	queryHandler = Handler{queryStore: store, osHandler: os}
}

func GetQueryHandler() *Handler {
	return &queryHandler
}

func (qh Handler) RunSavedQuery(name string) (ResultDTO, error) {
	logrus.Info("Running query: ", name)
	query, err := qh.fetchQuery(name)
	if err != nil {
		return ResultDTO{}, err
	}
	return qh.executeQuery(query)
}

func (qh Handler) RunCustomQuery(query string) (ResultDTO, error) {
	logrus.Info("Running query: ", query)
	return qh.executeQuery(query)
}

func (qh Handler) fetchQuery(name string) (string, error) {
	if query, ok := qh.queryStore.GetQuery(name); !ok {
		return "", fmt.Errorf("missing query")
	} else {
		return query, nil
	}
}

func (qh Handler) executeQuery(query string) (ResultDTO, error) {
	response, err := qh.osHandler.Execute(query)
	if err != nil {
		return ResultDTO{}, err
	}
	var result ResultDTO
	if err := result.UnmarshalArguments(cleanJSON(response)); err != nil {
		return ResultDTO{}, err
	}
	return result, nil
}

func cleanJSON(out bytes.Buffer) []byte {
	jsonRegex.Longest()
	cleanedJSON := jsonRegex.Find(out.Bytes())
	if !json.Valid(cleanedJSON) {
		logrus.Error(fmt.Errorf("invalid json"))
		return []byte{}
	}
	return cleanedJSON
}
