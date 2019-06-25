package rest

import (
	"log"
	"os"

	"github.com/AHerczeg/ostopus/head/api/rest"
	"github.com/AHerczeg/ostopus/head/api/rest/operation"
	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"
)

func StartServing(address string) {
	// TODO set address

	swaggerSpec, err := loads.Embedded(rest.SwaggerJSON, rest.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operation.NewOstopusAPI(swaggerSpec)
	server := rest.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "OStopus Head"
	parser.LongDescription = "The OStopus Head (server) API"

	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
