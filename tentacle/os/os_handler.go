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

func NewOSHandler() OSHandler {
	return stdOSHandler{}
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