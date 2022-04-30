#!/bin/bash
set -eo pipefail

#
# Created a .env file under the project root
# Below is an example
# 
# REDIS_CONNECTION=redis://@localhost:6379/1
# ALPHAVANTAGE_APIKEY=C227WD9W3LUVKVV9
# NDAYS=200
# SYMBOL=MSFT

REQUIRED_ENV=(
    ALPHAVANTAGE_APIKEY
    NDAYS
    SYMBOL
)

project_root="$(dirname "$0")/.."
cd "$project_root"

if [[ -f .env.default ]]; then
      export $(cat .env.default | xargs)
fi
if [[ -f .env ]]; then
      export $(cat .env | xargs)
fi

if [ -f .env ]
then
      export $(cat .env | xargs)
fi

for v in "${REQUIRED_ENV[@]}"; do
    if [[ -z "$(eval echo \$${v})" ]]; then
        echo "Please specify ${v} using environment variable" >&2
        exit 1
    fi
done

go run cmd/stock-price-grabber-server/main.go --port=3000
