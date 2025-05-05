#!/bin/bash
which venom
if [[ $? -ne 0 ]]; then
    ls ~/venom # checking if i was already installed from a previous run
    if [[ $? -ne 0 ]]; then
        echo "Venom is not installed. Installing venom binary."
        curl https://github.com/ovh/venom/releases/download/v1.1.0/venom.linux-amd64 -L -o ~/venom && chmod +x ~/venom
    fi
    VENOM=~/venom
else
    VENOM=$(which venom)
fi

E2E_TEST_MODE=1 ./dist/eigen_db&

$VENOM run e2e/ --output-dir=e2e/logs

E2E_EXIT_CODE=$?

pkill eigen_db
if [[ $? -ne 0 ]]; then
    echo "Failed to kill eigen_db process."
else
    echo "eigen_db process killed successfully."
fi

exit $E2E_EXIT_CODE