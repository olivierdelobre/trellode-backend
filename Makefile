BINARY_NAME=api
# if nothing is set, it will use the Dockerfile
compose_build: export DOCKERFILE = Dockerfile_localbuild
deploy_test: export DOCKERFILE = Dockerfile_localbuild
deploy_preprod: export DOCKERFILE = Dockerfile_localbuild

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

## build: build the binary (${BINARY_NAME})
build:
	CGO_ENABLED=0 GOOS=linux GOFLAGS="-ldflags=-s -ldflags=-w" go build -o server ./cmd/api/

## clean: clean go artefacts (binary included)
clean:
	go clean
	rm ${BINARY_NAME}

## compose_build: docker-compose build
compose_build: build
	docker-compose build

## compose_run: stop/build and relaunch docker-compose
compose_run: compose_stop compose_build
	docker-compose up

## compose_stop: docker-compose stop
compose_stop:
	docker-compose stop

## doc: makes documentation
.PHONY: doc
doc:
	swag init -g cmd/api/main.go

## local: launch docker-compose env
local: build compose_run

## release: test build and audit current code
release: test build audit

## build: buil the binary
run: build
	./${BINARY_NAME}

## test: launch quick tests
test: 
	go test -v ./...

## fulltest: launch full test on docker-compose environment
fulltest: 
	./run_tests.sh

## deploy_test: deploy in test
deploy_test: doc
	build_push.sh -i md-api-trellode -k ../md-api-infra/trellode -n md-api-test

## deploy_preprod: deploy in preprod
deploy_preprod: doc
	build_push.sh -i md-api-trellode -k ../md-api-infra/trellode -n md-api-preprod

## help: display this usage
.PHONY: help
help:
	@echo 'Usage:'
	@echo ${MAKEFILE_LIST}
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]
