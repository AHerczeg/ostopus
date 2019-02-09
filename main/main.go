package main

import (
	"context"
	"log"
	"os"
	"github.com/kolide/osquery-go"
)

func main() {
	server, err := osquery.NewExtensionManagerServer("foobar", os.Args[1])
	if err != nil {
		log.Fatalf("Error creating extension: %s\n", err)
	}

	// Create and register a new table plugin with the server.
	// table.NewPlugin requires the table plugin name,
	// a slice of Columns and a Generate function.
	server.RegisterPlugin(table.NewPlugin("foobar", FoobarColumns(), FoobarGenerate))
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}

	// FoobarColumns returns the columns that our table will return.
func FoobarColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("foo"),
		table.TextColumn("baz"),
	}
}

// FoobarGenerate will be called whenever the table is queried. It should return
// a full table scan.
func FoobarGenerate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	return []map[string]string{
		{
			"foo": "bar",
			"baz": "baz",
		},
		{
			"foo": "bar",
			"baz": "baz",
		},
	}, nil
}