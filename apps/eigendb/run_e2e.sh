#!/bin/bash
docker compose -f docker-compose-test.yml up -d --build

venom version
if $? -ne 0; then
    echo "Venom is not installed. Installing venom binary."
    curl https://github.com/ovh/venom/releases/download/v1.1.0/venom.linux-amd64 -L -o ./venom && chmod +x ./venom
    VENOM=./venom
else
    echo "Venom is already installed."
    VENOM=$(which venom)
fi

$VENOM run e2e/ --output-dir=e2e/logs
E2E_TEST_EXIT_CODE=$?

docker compose down

exit $E2E_TEST_EXIT_CODE