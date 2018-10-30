VERSION := $(shell git describe --abbrev=4 --dirty --always --tags)
Minversion := $(shell date)
BUILD_NODE_PAR = -ldflags "-X github.com/ioeXNetwork/ioeX.MainChain/config.Version=$(VERSION)" #-race

all:
	go build $(BUILD_NODE_PAR) -o ioex main.go

format:
	go fmt ./*

clean:
	rm -rf *.8 *.o *.out *.6
