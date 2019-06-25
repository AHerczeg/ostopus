# Basic go commands
GOCMD=go
GOBUILD=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
BINARY_NAME_HEAD=head-build
BINARY_NAME_TENTACLE=tentacled

.DEFAULT_GOAL := run.all

## source : http://www.oocities.org/spunk1111/aquatic.htm
define textlogo
                     ______
                 { /        \ }
                  / /o \  / o
                 |  \__/  \__/
                  \   ( ^ )  /           ___.--,
          _.._     \   uu   /     _.---'`__.-( (_.
   __.--'`_.. '.__.\    '--.\__.-' ,.--'`     `""`
  ( ,.--'`   ',__ /./;   ;, '.__.'`    __
  _`) )  .---.__.' / |   |\   \__..--""  """--.,_
 `---' .'.''-._.-'`_./  /\ '.  \ _.-~~~````~~~-._`-.__.'
       | |  .' _.-' |  |  \  \  '.               `~---`
        \ \/ .'     \  \   '. '-._)    ____   _____ _
         \/ /        \  \    `=.__`~-./ __ \ / ____| |
         / /\         `) )    / / `""| |\ | | (___ | |_ ___  _ __  _   _ ___
   , _.-'.'\ \        / /    ( (     | ) }| |\___ \| __/ _ \| '_ \| | | / __|
    `--~`   ) )    .-'.'      '.'.   |/ /_| |____) | || (_) | |_) | |_| \__ |
           (/`    ( (`          ) )  (/`___/|_____/ \__\___/| .__/ \__,_|___/
            `      '-;         (-'                          | |
                                                            |_|
endef


.PHONY: head
head:
	$(GOBUILD) -o $(BINARY_NAME_HEAD) -v ./head/main/

.PHONY: tentacle
tentacle:
	$(GOBUILD) -o $(BINARY_NAME_TENTACLE) -v ./tentacle/main/


.PHONY: swagger.head
swagger.head:
	swagger generate server --target ./head/api --name ostopus --spec ./head/api/swagger.yml --model-package model -a operation -s rest --exclude-main


.PHONY: build.head
build.head: test.head head

.PHONY: build.tentacle
build.tentacle: test.tentacle tentacle

.PHONY: run.head
run.head: build.head
	./$(BINARY_NAME_HEAD)

.PHONY: run.tentacle
run.tentacle: build.tentacle
	./$(BINARY_NAME_TENTACLE)

.PHONY: run.docker
run.docker: build.head
	docker build -t head-alpine -f head/Dockerfile .
	docker run -p 6060:6060 -p 7070:7070 head-alpine

.PHONY: run.all
run.all: logo run.docker run.tentacle

.PHONY: test.head
test.head:
	$(GOTEST) -v ./head/...

.PHONY: test.tentacle
test.tentacle:
	$(GOTEST) -v ./tentacle/...

.PHONY: test
test: test.head test.tentacle

.PHONY: logo
logo:
	$(info $(textlogo))

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f ./headd
	rm -f ./tentacled

.PHONY: sanitize
sanitize:
	go fmt ./head/...
	go fmt ./tentacle/...
    go vet -composites=false ./head/...
    go vet -composites=false ./tentacle/...


