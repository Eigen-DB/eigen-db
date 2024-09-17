#!/bin/sh

echo "Running tests..."
go test ./... -count=1 -v # running the tests. "-count=1" stops test caching
TEST_EXIT_CODE=$?
echo "Done."

exit $TEST_EXIT_CODE