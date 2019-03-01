package os

import (
	"bytes"
	"os/exec"
)

type OSHandler interface {
	Execute(string) (bytes.Buffer, error)
}

type stdOSHandler struct {
}

func (oh stdOSHandler) Execute(query string) (bytes.Buffer, error) {
	var out bytes.Buffer
	cmd := exec.Command("osqueryi", "--json", query)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return bytes.Buffer{}, err
	}
	return out, nil
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
