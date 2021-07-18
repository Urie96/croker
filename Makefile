.PHONE: all install

croker=$(shell which croker)
crokerd=$(shell which crokerd)

all: install

install: $(croker) $(crokerd)

$(croker): $(shell find croker -name \*.go)
	go get github.com/urie96/croker/croker

$(crokerd): $(shell find crokerd -name \*.go)
	go get github.com/urie96/croker/crokerd