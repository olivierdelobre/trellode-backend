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

The API should be available on port 8080 of current host

Get boards:
```
curl -v -X OPTIONS -H 'Authorization: Bearer 1' 'localhost:8080/trellode-api/v1/boards' | jq
curl -v -H 'Authorization: Bearer 1' 'localhost:8080/trellode-api/v1/boards' | jq
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

Get list members:
```
curl -v -H 'X-Krakend-UserType: service' -H 'X-Krakend-UserId: M00001' 'localhost:8080/trellode-api/v1/trellode/105179/members' | jq
curl -v -H 'X-Krakend-UserType: service' -H 'X-Krakend-UserId: M00001' 'localhost:8080/trellode-api/v1/trellode/105179/members?format=legacy' | jq
```

Search trellode:
```
curl -v -H 'X-Krakend-UserType: service' -H 'X-Krakend-UserId: M00001' 'localhost:8080/trellode-api/v1/trellode?type=batiments' | jq
curl -v -H 'X-Krakend-UserType: service' -H 'X-Krakend-UserId: M00001' 'localhost:8080/trellode-api/v1/trellode?type=batiments&format=legacy' | jq
```

Types:
```
curl -v -H 'X-Krakend-UserType: service' -H 'X-Krakend-UserId: M00001' 'localhost:8080/trellode-api/v1/types' | jq
curl -v -H 'X-Krakend-UserType: service' -H 'X-Krakend-UserId: M00001' 'localhost:8080/trellode-api/v1/types?format=legacy' | jq
```

Subtypes:
```
curl -v -H 'X-Krakend-UserType: service' -H 'X-Krakend-UserId: M00001' 'localhost:8080/trellode-api/v1/subtypes' | jq
curl -v -H 'X-Krakend-UserType: service' -H 'X-Krakend-UserId: M00001' 'localhost:8080/trellode-api/v1/subtypes?format=legacy' | jq
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
