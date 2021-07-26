.PHONE: all install

croker=$(shell which croker)

all: install

install: $(croker)

$(croker): $(shell find . -name \*.go)
	go get -ldflags="-w -s -X github.com/urie96/croker/version.Version=0.0.0" github.com/urie96/croker