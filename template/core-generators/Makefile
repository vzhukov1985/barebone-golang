GOPRIVATE=github.com/{{.orgName}}/{{.pkgRepoName}}
export

RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

install:
	@go install golang.org/x/tools/cmd/goimports@latest

gen-rest:
	@go run ./cmd/gen-rest/main.go $(RUN_ARGS)

gen-svc:
	@go run ./cmd/gen-svc/main.go $(RUN_ARGS)
	@cd ../core-$(word 2, $(MAKECMDGOALS)) && make mod-build

gen-ws:
	@go run ./cmd/gen-ws/main.go