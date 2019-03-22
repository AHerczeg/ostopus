# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
BINARY_NAME_HEAD=headd
BINARY_NAME_TENTACLE=tentacled

.DEFAULT_GOAL := run.all

## source : http://www.oocities.org/spunk1111/aquatic.htm
define textlogo
                     ______
                 { /        \ }    OStopus
                  / /o \  / o\
                 |  \__/  \__/
                  \   ( ^ )  /           ___.--,
          _.._     \   uu   /     _.---'`__.-( (_.
   __.--'`_.. '.__.\    '--. \_.-' ,.--'`     `""`
  ( ,.--'`   ',__ /./;   ;, '.__.'`    __
  _`) )  .---.__.' / |   |\   \__..--""  """--.,_
 `---' .'.''-._.-'`_./  /\ '.  \ _.-~~~````~~~-._`-.__.'
       | |  .' _.-' |  |  \  \  '.               `~---`
        \ \/ .'     \  \   '. '-._)
         \/ /        \  \    `=.__`~-.
         / /\         `) )    / / `"".`\
   , _.-'.'\ \        / /    ( (     / /
    `--~`   ) )    .-'.'      '.'.  | (
           (/`    ( (`          ) )  '-;
            `      '-;         (-'
endef

head:
	$(GOBUILD) -o ostopus/$(BINARY_NAME_HEAD) -v ./head/main/

tentacle:
	$(GOBUILD) -o ostopus/$(BINARY_NAME_TENTACLE) -v ./tentacle/main/

build.head: test.head head

build.tentacle: test.tentacle tentacle

run.head: build.head
	./$(BINARY_NAME_HEAD)

run.tentacle: build.tentacle
	./$(BINARY_NAME_TENTACLE)

run.docker: build.head
	docker build -t head-alpine -f head/Dockerfile .

run.all: run.docker run.tentacle

test.head:
	$(GOBUILD) -v ./head/...

test.tentacle:
	$(GOBUILD) -v ./tentacle/...

sanitize:
	go fmt ./head/...
	go fmt ./tentacle/...
    go vet -composites=false ./head/...
    go vet -composites=false ./tentacle/...

clean:
	$(GOCLEAN)
	rm -f ./headd
	rm -f ./tentacled