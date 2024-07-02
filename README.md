# Trellode API

REST API for trellode


## Overview 

The module will provide :
* A REST API

## Operations

### Make usage

To ease use, make commands are made available:

* make build (build api)
* make test (launch unit tests)
* make fulltest (launch unit tests+E2E tests)
* make local (build and launch local containers)
* make release
* ...

Try ``make help`` for more details


### Docker operations

#### Start docker containers

```
docker-compose up -d
```

The API should be available on port 8080 of current host.

Register:
```
# Password too short
curl -v -X POST -H 'Content-Type: application/json' -d '{"email":"o.livier@gmail.com", "firstname":"Olivier", "lastname": "Delobre", "password": "azerty"}' 'localhost:8080/trellode-api/v1/users/register' | jq

# Bad password strength
curl -v -X POST -H 'Content-Type: application/json' -d '{"email":"o.livier@gmail.com", "firstname":"Olivier", "lastname": "Delobre", "password": "azertyazerty"}' 'localhost:8080/trellode-api/v1/users/register' | jq

# Alrighty
curl -v -X POST -H 'Content-Type: application/json' -d '{"email":"o.livier@gmail.com", "firstname":"Olivier", "lastname": "Delobre", "password": "azerTY1234+"}' 'localhost:8080/trellode-api/v1/users/register' | jq
```

Authenticate:
```
# Bad password
curl -v -X POST -H 'Content-Type: application/json' -d '{"email":"o.livier@gmail.com", "password": "azerty"}' 'localhost:8080/trellode-api/v1/users/authenticate' | jq

# Alrighty
curl -v -X POST -H 'Content-Type: application/json' -d '{"email":"o.livier@gmail.com", "password": "azerTY1234+"}' 'localhost:8080/trellode-api/v1/users/authenticate' | jq
```

Get boards:
```
curl -v -X OPTIONS -H 'Authorization: Bearer 1' 'localhost:8080/trellode-api/v1/boards' | jq
curl -v -H 'Authorization: Bearer 1' 'localhost:8080/trellode-api/v1/boards' | jq
curl -v -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjE1YmExZjRlLWQ2ODItNGY1NS04NTgzLWRkMzEwYTY3MjVlNyIsImVtYWlsIjoiby5saXZpZXJAZ21haWwuY29tIiwiZmlyc3RuYW1lIjoiT2xpdmllciIsImxhc3RuYW1lIjoiRGVsb2JyZSIsInByb2ZpbGUiOiJ1c2VyIiwiZXhwIjoxNzE5OTM2ODE2LCJpYXQiOjE3MTk5MjYwMTZ9.O427I0y2QVxBXc9r4yq3KCkw_vAHwRQzHCcz-6Y8akI' 'localhost:8080/trellode-api/v1/boards' | jq
```

Get board:
```
curl -v -X OPTIONS -H 'Authorization: Bearer 1' 'localhost:8080/trellode-api/v1/boards/1' | jq
curl -v -H 'Authorization: Bearer 1' 'localhost:8080/trellode-api/v1/boards/1' | jq
```

Create board:
```
curl -v -H 'Authorization: Bearer 1' -d '{"title":"board theta"}' 'localhost:8080/trellode-api/v1/boards' | jq
```

Create background:
```
curl -v -X POST -H 'Authorization: Bearer 1' -H 'Content-Type: application/json' -d '{"data":"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADIAAAAyCAIAAACRXR/mAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAyJpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADw/eHBhY2tldCBiZWdpbj0i77u/IiBpZD0iVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkIj8+IDx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IkFkb2JlIFhNUCBDb3JlIDUuMC1jMDYwIDYxLjEzNDc3NywgMjAxMC8wMi8xMi0xNzozMjowMCAgICAgICAgIj4gPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4gPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIgeG1sbnM6eG1wPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvIiB4bWxuczp4bXBNTT0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL21tLyIgeG1sbnM6c3RSZWY9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZVJlZiMiIHhtcDpDcmVhdG9yVG9vbD0iQWRvYmUgUGhvdG9zaG9wIENTNSBNYWNpbnRvc2giIHhtcE1NOkluc3RhbmNlSUQ9InhtcC5paWQ6RDUxRjY0ODgyQTkxMTFFMjk0RkU5NjI5MEVDQTI2QzUiIHhtcE1NOkRvY3VtZW50SUQ9InhtcC5kaWQ6RDUxRjY0ODkyQTkxMTFFMjk0RkU5NjI5MEVDQTI2QzUiPiA8eG1wTU06RGVyaXZlZEZyb20gc3RSZWY6aW5zdGFuY2VJRD0ieG1wLmlpZDpENTFGNjQ4NjJBOTExMUUyOTRGRTk2MjkwRUNBMjZDNSIgc3RSZWY6ZG9jdW1lbnRJRD0ieG1wLmRpZDpENTFGNjQ4NzJBOTExMUUyOTRGRTk2MjkwRUNBMjZDNSIvPiA8L3JkZjpEZXNjcmlwdGlvbj4gPC9yZGY6UkRGPiA8L3g6eG1wbWV0YT4gPD94cGFja2V0IGVuZD0iciI/PuT868wAAABESURBVHja7M4xEQAwDAOxuPw5uwi6ZeigB/CntJ2lkmytznwZFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYW1qsrwABYuwNkimqm3gAAAABJRU5ErkJggg=="}' 'localhost:8080/trellode-api/v1/backgrounds' | jq
```

Get background:
```
curl -v -H 'Authorization: Bearer 1' 'localhost:8080/trellode-api/v1/backgrounds/1' | jq
```

Get backgrounds:
```
curl -v -H 'Authorization: Bearer 1' 'localhost:8080/trellode-api/v1/backgrounds' | jq
```

Healthcheck
```
curl -v 'localhost:8080/healthcheck'
```

Liveness
```
curl -v 'localhost:8080/liveness'
curl -v 'localhost:8080/liveness?format=json' | jq
curl -v 'localhost:8080/liveness?format=metrics'
```

### Binaries operations
Each binary requires its configuration (database credentials...) to be available through a .env file located in the same directory as the binary.
You can use env.sample as template
```
cp env.sample .env
# Modify .env to update database credentials
```

Also note that the binary has to run on an host which is allowed to access the target database (ex.: test-cadibtch if you want to connect to test-cadidb)

#### Build all binaries

```
go build  -o . ./...
```

#### Build API server binary

```
go build cmd/api/
```

#### Launch API binary

```
./api # or ./api &
```

## Tests
### Services unit tests
```
go test -v internal/room/service_test.go
```
### API test
Run locally
```
go test -v internal/api/api_test.go
```
Run in container
```
docker-compose up -d
docker exec trellode-api-test bash -c "go test -v internal/api/api_test.go"
```

## Generate mocks with mockery
See https://medium.com/deliveryherotechhub/mocking-an-interface-using-mockery-in-go-afbcb83cc773

### Install mockery
```
brew install mockery
```
### Generate a mock
```
mockery --dir=internal/list --name=ListRepositoryInterface --filename=ListRepositoryInterface.go --output=internal/mocks/repomocks --outpkg=repomocks
mockery --dir=internal/listtype --name=TypeRepositoryInterface --filename=TypeRepositoryInterface.go --output=internal/mocks/repomocks --outpkg=repomocks
mockery --dir=internal/listsubtype --name=SubtypeRepositoryInterface --filename=SubtypeRepositoryInterface.go --output=internal/mocks/repomocks --outpkg=repomocks
```

## Deploy on OpenShift
You need the generic build_push.sh script from md-tools repo:
test:
```
build_push.sh -i md-api-trellode -k ../md-api-infra/trellode -n md-api-test
```
prod:
```
build_push.sh -t <tag name> -i md-api-trellode -k ../md-api-infra/trellode -n md-api-prod
```

## Generate documentation
Run following command int the project's root to generate Swagger documentation:
```
swag init -d cmd/api,internal/api --pd
```
This will create a *docs* folder at the root of the project.

Then go to http://localhost:8080/trellode-api/v1/docs/index.html to see generated documentation.
