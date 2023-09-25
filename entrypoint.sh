#!/bin/sh
set -eo pipefail

reflex -g 'hapi-*-ssm-secrets/*' --shutdown-timeout=30000ms -s ./happy-api "$@"
