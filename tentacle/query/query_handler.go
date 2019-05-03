package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"ostopus/tentacle/os"
	"regexp"
)

var (
	queryHandler StdHandler
	jsonRegex     = regexp.MustCompile("{\\s*(\"[^\"]*\":\"[^\"]*\"\\s*)+(,\\s*(\"[^\"]*\":\"[^\"]*\")\\s*)*}")
)

type Handler interface {
	RunSavedQuery(string) (ResultDTO, error)
	RunCustomQuery(string) (ResultDTO, error)
	fetchQuery(string) (string, error)
	executeQuery(string) (ResultDTO, error)
}

type StdHandler struct {
	queryStore Store
	osHandler  os.Handler
}

func InitQueryHandler(store Store, os os.Handler) {
	queryHandler = StdHandler{queryStore: store, osHandler: os}
}

func GetQueryHandler() *StdHandler {
	return &queryHandler
}

func (qh StdHandler) RunSavedQuery(name string) (ResultDTO, error) {
	logrus.Info("Running query: ", name)
	query, err := qh.fetchQuery(name)
	if err != nil {
		return ResultDTO{}, err
	}
	return qh.executeQuery(query)
}

func (qh StdHandler) RunCustomQuery(query string) (ResultDTO, error) {
	logrus.Info("Running query: ", query)
	return qh.executeQuery(query)
}

func (qh StdHandler) fetchQuery(name string) (string, error) {
	if query, ok := qh.queryStore.GetQuery(name); !ok {
		return "", fmt.Errorf("missing query")
	} else {
		return query, nil
	}
}

func (qh StdHandler) executeQuery(query string) (ResultDTO, error) {
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


/** Mocks **/

type MockQueryHandler struct {
	mock.Mock
}

func (mq *MockQueryHandler) RunSavedQuery(name string) (ResultDTO, error) {
	args := mq.Called(name)
	return args.Get(0).(ResultDTO), args.Error(1)
}

func (mq *MockQueryHandler) RunCustomQuery(query string) (ResultDTO, error) {
	args := mq.Called(query)
	return args.Get(0).(ResultDTO), args.Error(1)
}

func (mq *MockQueryHandler) fetchQuery(name string) (string, error) {
	args := mq.Called(name)
	return args.String(0), args.Error(1)
}

func (mq *MockQueryHandler) executeQuery(query string) (ResultDTO, error) {
	args := mq.Called(query)
	return args.Get(0).(ResultDTO), args.Error(1)
}