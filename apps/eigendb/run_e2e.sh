#!/bin/bash
docker compose -f docker-compose-test.yml up -d --build

venom run e2e/ --output-dir=e2e/logs
E2E_TEST_EXIT_CODE=$?

docker compose down

exit $E2E_TEST_EXIT_CODE