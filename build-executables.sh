#!/usr/bin/env bash
# -------------------------
# Use release.sh to release 
# -------------------------
# - it does stuff like this -
# git tag -a v0.1.0 -m "Release info"
# git push origin v0.1.0
# goreleaser --rm-dist

# To see all supported buids
# > go tool dist list

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi

platforms=("windows/amd64" "windows/386" "darwin/amd64", "linux/amd64")

for platform in "${platforms[@]}"
do
  platform_split=(${platform//\// })
  GOOS=${platform_split[0]}
  GOARCH=${platform_split[1]}
  output_name='tmint-'$GOOS'-'$GOARCH
  if [ $GOOS = "windows" ]; then
    output_name+='.exe'
  fi

  env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package
  if [ $? -ne 0 ]; then
    echo 'An error has occurred! Aborting the script execution...'
    exit 1
  fi
done
