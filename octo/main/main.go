package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"ostopus/octo/query"
)

type result struct {
	Arguments string `json:"arguments"`
	Device    string `json:"device"`
	Path      string `json:"path"`
	Version   string `json:"version"`
}

func main() {
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
}
