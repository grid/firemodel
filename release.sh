#!/bin/bash

set -e

if [[ ! -z "${DEBUG}" ]]; then
  set -x
fi

echo "Checking for dirty working directory…"
if [[ ! -z "$(git diff --exit-code)" ]]; then
  echo error: Dirty working directory.
  exit 1
fi

if [[ ! -z "$(git diff --cached --exit-code)" ]]; then
  echo error: Dirty index.
  exit 1
fi

echo "Preflight…"

go test ./... || (echo "error: Tests failed."; exit 1)

while [[ "$(hub ci-status)" == "pending" ]]; do
  echo "Waiting for CI…";
  sleep 1
done

hub ci-status || (echo "error: CI failed"; exit 1)

echo Fetching latest tags…
git fetch --tags origin

echo "10 most recent releases:"
hub release --include-drafts --limit=10

if [[ -z "${NEW_TAG}" ]]; then
  read -r -p "Enter tag for new release: " NEW_TAG
fi

echo "${NEW_TAG?}"

git tag -am "${NEW_TAG?}" "${NEW_TAG?}"

echo "Building firemodel…"

mkdir -p .build

GOARCH=amd64

for goos in "darwin" "linux" "windows"; do
  GOOS=${goos?} \
    go get ./...

  GOOS=${goos?} \
    go build \
    -ldflags "-X github.com/visor-tax/firemodel/version.Version=${NEW_TAG?}" \
    -o "./.build/firemodel-${goos?}-${GOARCH?}" ./firemodel/main.go
done

read -r -p "Release on GitHub? [yes/no] " PUSH_TAG

case "${PUSH_TAG?}" in
  y | yes | YES)
    echo "Pushing tag…"
    git push origin "${NEW_TAG?}"

    echo "Creating draft release…"

    hub release \
      create \
      --message="${NEW_TAG?}" \
      --attach=.build/firemodel-darwin-amd64 \
      --attach=.build/firemodel-linux-amd64 \
      --attach=.build/firemodel-windows-amd64 \
      --browse \
      "${NEW_TAG?}"

    say "Released firemodel ${NEW_TAG?}"
    ;;
  *)
    echo "Skipping tag push."
    ;;
esac
