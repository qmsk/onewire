export GOBIN=$(CURDIR)/bin

PREFIX=/opt/qmsk-onewire

bin:
	go install -v github.com/qmsk/onewire/cmd/...

install:
	install -d ${PREFIX}
	install -d ${PREFIX}/bin
	install -t ${PREFIX}/bin/ -m 0755 bin/*
