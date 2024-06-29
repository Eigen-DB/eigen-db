#!/bin/sh

cp ./lib/libhnsw.so /usr/local/lib/ # copying shared library into /usr/local/lib for tests to run
echo "Running tests..."
go test ./... -count=1 # running the tests. "-count=1" stops test caching
TEST_EXIT_CODE=$?

# cleanup
rm -rf /usr/local/lib/libhnsw.so
echo "Done."

exit $TEST_EXIT_CODE