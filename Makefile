BINARY=debezium_test
PROJECT_NAME=debezium_test
ENV ?= dev
LOG_LEVEL ?= debug
IMAGE_NAME ?= debezium_test
IMAGE_TAG ?= latest

all: clean build

clean:
	@echo "--> Target directory clean up"
	rm -rf ./.build/target
	rm -f ${BINARY}

lint:
	@echo "--> Running Golang linter (unused variable / function warning are skipped)"
	golangci-lint run --exclude 'unused' --build-tags integ

test-unit:
	@echo "--> Testing: UNIT tests (no cache)"
	go test -race -v $(shell go list ./... | grep -v 'tests/') -count=1

test-unit-ci:
	gotestsum --junitfile test_reports/unit-tests.xml -- -race $(shell go list ./... | grep -v 'tests/') -count=1

test-integ:
	@echo "--> Testing: Integration tests (no cache)"
	DEBEZIUM_TEST_DB_HOST=127.0.0.1 \
    DEBEZIUM_TEST_DB_PORT=5432 \
    DEBEZIUM_TEST_DB_USER=postgres \
    DEBEZIUM_TEST_DB_PASSWORD=1q2w3e \
    DEBEZIUM_TEST_DB_NAME=debezium_test \
    PUBSUB_EMULATOR_HOST=127.0.0.1:8681 GCLOUD_PROJECT=local-docker GOOGLE_CLOUD_PROJECT=local-docker go test --tags=integ -v -p 1 ./tests/... -count=1

test-integ-ci:
	gotestsum --junitfile test_reports/integ-tests.xml -- --tags=integ -v ./tests/... -count=1

build:
	go build -o ${BINARY} ./cmd

build-docker:
	./.build/build_docker.sh ./cmd ${BINARY} ${IMAGE_NAME} ${IMAGE_TAG}

run-compose:
	docker-compose --project-name ${PROJECT_NAME} -f docker-compose.yml down
	docker-compose --project-name ${PROJECT_NAME} -f docker-compose.yml up --force-recreate --remove-orphans --renew-anon-volumes