#!/bin/bash
set -eo pipefail

REQUIRED_ENV=(
    SYMBOL
    NDAYS
    KUBE_CONTEXT
    ALPHAVANTAGE_APIKEY_BASE64
    MIN_REPLICA
    MAX_REPLICA
    CPU_THRESHOLD
    MEM_THRESHOLD
    IMAGE_REPO
    IMAGE_TAG
    CPU_LIMIT
    MEM_LIMIT
    CPU_REQUEST
    MEM_REQUEST
)

if [[ -z "$ENV" ]]; then
    echo Please specify the environment using ENV variable >&2
    exit 1
fi

project_root="$(dirname "$0")/.."
cd "$project_root"

if [[ -f .default.env ]]; then
      export $(cat .default.env | xargs)
fi
if [ -f ".env.${ENV}" ]; then
    export $(cat ".env.${ENV}" | xargs)
fi

for v in "${REQUIRED_ENV[@]}"; do
    if [[ -z "$(eval echo \$${v})" ]]; then
        echo "Please specify ${v} using environment variable" >&2
        exit 1
    fi
done

for f in deploy/*.tmpl; do
    cat "$f" | envsubst | kubectl --context "$KUBE_CONTEXT" apply -f -
done
