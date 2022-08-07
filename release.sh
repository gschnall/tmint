#!/usr/bin/env bash


# Example
# -------
# ./release.sh 0.3.1 'Fix some bug'

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

function check_for_git_changes {
  if [[ $(git diff --stat) != '' ]]; then
    echo 'There are uncommitted changes in your git directory'
    echo "-You'll need to commit, or remove, them before goreleaser will work."
    exit 1
  fi

  local untracked_files=$(git status --porcelain 2>/dev/null| grep "^??" | wc -l)
  if [[ $untracked_files > 0 ]]; then
    echo "You have untracked_files in the repo."
    echo "-You'll need to commit, or remove, them before goreleaser will work."
    exit 1
  fi
}

function main {
  check_for_git_changes

  echo " Release"
  echo "---------"
  echo "version: $version"
  echo "message: $message"

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
