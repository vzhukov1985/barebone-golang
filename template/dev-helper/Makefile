GOPRIVATE=github.com/{{.orgName}}/{{.pkgRepoName}}
export

include .envs/.env.development.local

RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

mod-upd:
	@echo Syncing and updating packages...
	-@rm go.mod
	-@rm go.sum
	@go mod init dev-helper
	@go mod tidy

build:
	go build -o ./bin/dev-helper ./cmd/main.go

run:
	go run cmd/main.go $(RUN_ARGS)