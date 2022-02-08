APPNAME=gnote
GIT_COMMIT:=$(shell git rev-list -1 HEAD)
GIT_TAG:=$(shell git describe --tags $(git rev-list --tags --max-count=1))
LDFLAGS:=-X main.GitCommit=${GIT_COMMIT} -X main.Version=${GIT_TAG} -s -w

build:
	go build -o "$(APPNAME)" -ldflags "$(LDFLAGS)" ./cmd/gnote/...

dev:
	go build -o "$(APPNAME)" ./cmd/gnote/...

install:
	go install -o "$(APPNAME)" -ldflags "$(LDFLAGS)" ./cmd/gnote/...

test:
	go test -v ./...

clean:
	rm "$(APPNAME)"
