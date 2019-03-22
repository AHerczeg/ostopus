# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
BINARY_NAME_HEAD=headd
BINARY_NAME_TENTACLE=tentacled

head:
	$(GOBUILD) -o $(BINARY_NAME_HEAD) -v ./head/main/

tentacle:
	$(GOBUILD) -o $(BINARY_NAME_TENTACLE) -v ./tentacle/main/

build.head:
	test.head
	head

build.tentacle:
	test.head
	tentacle

run.head:
	build.head
	./$(BINARY_NAME_HEAD)

run.tentacle:
	build.tentacle
	./$(BINARY_NAME_TENTACLE)

test.head:
	$(GOBUILD) -v ./head/...

test.tentacle:
	$(GOBUILD) -v ./tentacle/...

clean:
	$(GOCLEAN)
	rm -f ./headd
	rm -f ./tentacled