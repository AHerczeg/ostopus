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

/*

import "bytes"

fmt.Println("Query: kernel_info")
query, ok := query.GetQuery("kernel_info")
if !ok {
fmt.Println("error: unexpected query")
}
cmd := exec.Command("osqueryi", "--json", query)
var out bytes.Buffer
cmd.Stdout = &out
err := cmd.Run()
if err != nil {
log.Fatal(err)
}
var res result
bytes, err := out.ReadBytes(']')
if err != nil {
log.Fatal(err)
}
json.Unmarshal(bytes[1:len(bytes)-1], &res)
fmt.Printf("arguments:\t | %s\n", res.Arguments)
fmt.Printf("device:\t\t | %s\n", res.Device)
fmt.Printf("path:\t\t | %s\n", res.Path)
fmt.Printf("version:\t | %s\n", res.Version)
*/
