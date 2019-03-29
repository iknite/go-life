.PHONY: prepare install release

ARTIFACTS_DIR=artifacts/${VERSION}
GITHUB_USERNAME=sachaos

test:
	go test -v ./format/...

prepare:
	go mod download
	go generate github.com/iknite/raft-life/preset

install: prepare
	go install

release: prepare
	GOOS=windows GOARCH=amd64 go build -o $(ARTIFACTS_DIR)/go-life_windows_amd64
	GOOS=darwin GOARCH=amd64 go build -o $(ARTIFACTS_DIR)/go-life_darwin_amd64
	GOOS=linux GOARCH=amd64 go build -o $(ARTIFACTS_DIR)/go-life_linux_amd64
	ghr -u $(GITHUB_USERNAME) -t $(shell cat github_token) --replace ${VERSION} $(ARTIFACTS_DIR)
