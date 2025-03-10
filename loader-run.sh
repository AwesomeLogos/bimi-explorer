#!/usr/bin/env bash
#
# run locally
#

set -o errexit
set -o pipefail
set -o nounset

if [ -f ".env" ]; then
	echo "INFO: loading .env file"
	export $(cat .env)
else
	echo "INFO: .env file not found"
fi

go run cmd/loader/bulkLoader.go $@
