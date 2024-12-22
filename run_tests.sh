#!/bin/bash

if [[ $WORKFLOW == 1 ]]
then
    DOCKER_COMPOSE=./docker-compose
    VENOM=./venom
    BUILD=1
else
    DOCKER_COMPOSE=docker-compose
    VENOM=venom
fi

echo "UNIT tests:"
go test ./... -count=1 -v # running the tests. "-count=1" stops test caching
UNIT_TEST_EXIT_CODE=$?

echo "END-TO-END tests:"
if [[ $BUILD == 1 ]]
then
    E2E_TEST_MODE=1 $DOCKER_COMPOSE -f docker-compose-test.yml up -d --build
else
    E2E_TEST_MODE=1 $DOCKER_COMPOSE -f docker-compose-test.yml up -d
fi
$VENOM run e2e/ --output-dir=e2e/logs
E2E_TEST_EXIT_CODE=$?

$DOCKER_COMPOSE down
echo "Done."

if [[ $E2E_TEST_EXIT_CODE != 0 || $UNIT_TEST_EXIT_CODE != 0 ]]
then
    exit 1
fi

exit 0 