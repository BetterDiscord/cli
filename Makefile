.PHONY: remove format

export GOFLAGS=-ldflags=-s -trimpath

DIR?=build

$(DIR):
	mkdir -p $(DIR)

$(DIR)/bdcli_darwin_amd64: $(DIR)
	env GOOS=darwin GOARCH=amd64 go build -o $@

$(DIR)/bdcli_darwin_arm64: $(DIR)
	env GOOS=darwin GOARCH=arm64 go build -o $@

$(DIR)/bdcli_linux_386: $(DIR)
	env GOOS=linux GOARCH=386 go build -o $@

$(DIR)/bdcli_linux_amd64: $(DIR)
	env GOOS=linux GOARCH=amd64 go build -o $@

$(DIR)/bdcli_linux_arm64: $(DIR)
	env GOOS=linux GOARCH=arm64 go build -o $@

$(DIR)/bdcli_linux_arm: $(DIR)
	env GOOS=linux GOARCH=arm go build -o $@

$(DIR)/bdcli_windows_386.exe: $(DIR)
	env GOOS=windows GOARCH=386 go build -o $@

$(DIR)/bdcli_windows_amd64.exe: $(DIR)
	env GOOS=windows GOARCH=amd64 go build -o $@

remove:
	@rm -rf builds

format:
	@go fmt ./...