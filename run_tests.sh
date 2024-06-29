#!/bin/sh

echo "Running tests..."
CGO_LDFLAGS="-L./lib -lhnsw" go test ./... -count=1 # running the tests. "-count=1" stops test caching
TEST_EXIT_CODE=$?
echo "Done."

exit $TEST_EXIT_CODE