#!/bin/bash

TIMESTAMP=$(date +%s%N)
mkdir -p e2e/logs

which venom
if [[ $? -ne 0 ]]; then
    ls ~/venom # checking if i was already installed from a previous run
    if [[ $? -ne 0 ]]; then
        echo "Venom is not installed. Installing venom binary."
        curl https://github.com/ovh/venom/releases/download/v1.2.0/venom.linux-amd64 -L -o ~/venom && chmod +x ~/venom
    fi
    VENOM=~/venom
else
    VENOM=$(which venom)
fi

E2E_TEST_MODE=1 ./dist/eigen_db > e2e/logs/eigen_db_$TIMESTAMP.log 2>&1 &
if [[ $? -ne 0 ]]; then
    echo "Failed to start eigen_db process."
    exit 1
else
    echo "eigen_db process started successfully. Runing E2E tests..."
fi

sleep 2 # hacky way to lower chances of race condition related flakiness

$VENOM run e2e/ --output-dir=e2e/logs

E2E_EXIT_CODE=$?

pkill eigen_db
if [[ $? -ne 0 ]]; then
    echo "Failed to kill eigen_db process."
else
    echo "eigen_db process killed successfully."
fi

echo "A log of eigen_db's output can be found in e2e/logs/eigen_db_$TIMESTAMP.log"

# cleanup
rm -rf eigen/*

exit $E2E_EXIT_CODE