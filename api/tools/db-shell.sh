#!/bin/bash
set -euox pipefail

docker compose run -it db psql -U postgres -h db
