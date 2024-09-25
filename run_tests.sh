#!/bin/bash

if [[ $WORKFLOW == 1 ]]
then
    DOCKER_COMPOSE=./docker-compose
    VENOM=./venom
else
    DOCKER_COMPOSE=docker-compose
    VENOM=venom
fi

echo "UNIT tests:"
go test ./... -count=1 -v # running the tests. "-count=1" stops test caching
TEST_EXIT_CODE=$?
if [[ $TEST_EXIT_CODE != 0 ]]
then
    exit $TEST_EXIT_CODE
fi

echo "INTEGRATION tests:"
if [[ $BUILD == 1 ]]
then
    TEST_MODE=1 $DOCKER_COMPOSE -f docker-compose-test.yml up -d --build
else
    TEST_MODE=1 $DOCKER_COMPOSE -f docker-compose-test.yml up -d
fi
$VENOM run integration_tests/ --output-dir=integration_tests/logs
TEST_EXIT_CODE=$?

$DOCKER_COMPOSE down
echo "Done."

exit $TEST_EXIT_CODE