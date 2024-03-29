export GOPRIVATE=github.com/ecumenos
export SHELL=/bin/sh
export GOEXPERIMENT := loopvar

.PHONY: all
all: tidy check fmt lint test mock tidy openapi tidy

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test-short
test-short: ## Run tests, skipping slower integration tests
	go test -test.short ./...

.PHONY: test-interop
test-interop: ## Run tests, including local interop (requires services running)
	go clean -testcache && go test -tags=localinterop ./...

.PHONY: coverage-html
coverage-html: ## Generate test coverage report and open in browser
	go test ./... -coverpkg=./... -coverprofile=test-coverage.out
	go tool cover -html=test-coverage.out

.PHONY: lint
lint: ## Verify code style and run static checks
	go vet -asmdecl -assign -atomic -bools -buildtag -cgocall -copylocks -httpresponse -loopclosure -lostcancel -nilfunc -printf -shift -stdmethods -structtag -tests -unmarshal -unreachable -unsafeptr -unusedresult ./...
	test -z $(gofmt -l ./...)

.PHONY: fmt
fmt: ## Run syntax re-formatting (modify in place)
	go fmt ./...

.PHONY: check
check: ## Compile everything, checking syntax (does not output binaries)
	go build ./...

.PHONY: mock
mock:
	sh scripts/mocks.sh

.PHONY: openapi
openapi:
	sh scripts/docs.sh

.env:
	if [ ! -f ".env" ]; then cp example.dev.env .env; fi

# PDS
.PHONY: run-dev-pds
run-dev-pds: .env ## Runs pds for local dev
	export API_LOCAL=true && go run cmd/pds/*.go run-api-server

.PHONY: run-dev-pds-admin
run-dev-pds-admin: .env ## Runs pds for local admin
	export API_LOCAL=true && go run cmd/pds/*.go run-admin-server

.PHONY: migrate-up-pds
migrate-up-pds: .env
	export API_LOCAL=true && go run cmd/pds/*.go migrate-up

.PHONY: migrate-down-pds
migrate-down-pds: .env
	export API_LOCAL=true && go run cmd/pds/*.go migrate-down

.PHONY: build-pds-image
build-pds-image:
	docker build -t pds -f cmd/pds/Dockerfile .

.PHONY: run-pds-image
run-pds-image:
	docker run -p 9090:9090 pds /pds  run-api-server

# Orbis Socius
.PHONY: run-dev-orbis-socius
run-dev-orbis-socius: .env ## Runs orbis socius for local dev
	export API_LOCAL=true && go run cmd/orbissocius/*.go run-api-server

.PHONY: run-dev-orbis-socius-admin
run-dev-orbis-socius-admin: .env ## Runs orbis socius for local admin
	export API_LOCAL=true && go run cmd/orbissocius/*.go run-admin-server

.PHONY: migrate-up-orbis-socius
migrate-up-orbis-socius: .env
	export API_LOCAL=true && go run cmd/orbissocius/*.go migrate-up

.PHONY: migrate-down-orbis-socius
migrate-down-orbis-socius: .env
	export API_LOCAL=true && go run cmd/orbissocius/*.go migrate-down

.PHONY: build-orbis-socius-image
build-orbis-socius-image:
	docker build -t orbissocius -f cmd/orbissocius/Dockerfile .

.PHONY: run-orbis-socius-image
run-orbis-socius-image:
	docker run -p 9091:9091 orbissocius /orbissocius  run-api-server

# Zookeeper
.PHONY: run-dev-zookeeper
run-dev-zookeeper: .env ## Runs zookeeper for local dev
	export API_LOCAL=true && go run cmd/zookeeper/*.go run-api-server

.PHONY: run-dev-zookeeper-admin
run-dev-zookeeper-admin: .env ## Runs zookeeper-admin for local dev
	export API_LOCAL=true && go run cmd/zookeeper/*.go run-admin-server

.PHONY: migrate-up-zookeeper
migrate-up-zookeeper: .env
	export API_LOCAL=true && go run cmd/zookeeper/*.go migrate-up

.PHONY: migrate-down-zookeeper
migrate-down-zookeeper: .env
	export API_LOCAL=true && go run cmd/zookeeper/*.go migrate-down

.PHONY: seeds-zookeeper
seeds-zookeeper: .env
	export API_LOCAL=true && go run cmd/zookeeper/*.go run-seeds

.PHONY: build-zookeeper-image
build-zookeeper-image:
	docker build -t zookeeper -f cmd/zookeeper/Dockerfile .

.PHONY: run-zookeeper-image
run-zookeeper-image:
	docker run -p 9092:9092 zookeeper /zookeeper  run-api-server
