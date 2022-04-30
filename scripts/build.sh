#!/bin/bash
set -eo pipefail

REQUIRED_ENV=(
    IMAGE_REPO
    IMAGE_TAG
)

project_root="$(dirname "$0")/.."
cd "$project_root"

if [[ -f .env.default ]]; then
      export $(cat .env.default | xargs)
fi
if [[ -f .env ]]; then
      export $(cat .env | xargs)
fi

for v in "${REQUIRED_ENV[@]}"; do
    if [[ -z "$(eval echo \$${v})" ]]; then
        echo "Please specify ${v} using environment variable" >&2
        exit 1
    fi
done

docker build -f build/Dockerfile -t "$IMAGE_REPO:$IMAGE_TAG" . 
docker push "$IMAGE_REPO:$IMAGE_TAG"
