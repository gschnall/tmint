#!/usr/bin/env bash

tag=$1
message="$2"
version="v$1"

if [[ -z "$tag" ]]; then
  latestTag=$(bash ./get-latest-tag.sh)
  echo "usage:   $0 <tag> <release-message>"
  echo "example: $0 0.1.0 'Message about release'"
  echo -e "\n $latestTag <-- The last created tag"
  exit 1
elif [[ -z "$message" ]]; then
  latestTag=$(bash ./get-latest-tag.sh)
  echo "usage:   $0 <tag> <release-message>"
  echo "example: $0 0.1.0 'Message about release'"
  echo ""
  echo -e "\n $latestTag <-- The last created tag"
  exit 1
fi

echo "version: $version"
echo "message: $message"

function main {
  git tag -a "$version" -m "$message"&&git push origin "$version"&&goreleaser --rm-dist

  echo "____________________"
  echo "|"
  echo "| version: $version"
  echo "| message: $message"
  echo "|"
  echo "| Released"
  echo "-------------------"
}
main

