#!/bin/bash -ex
git config --global user.email "infra-eng+gha@chanzuckerberg.com"
git config --global user.name "GitHub Actions Bot"

awk -i inplace 'BEGIN { FS = "." } ; { print $1 "." $2 "." ++$3 }' VERSION
version=$(cat VERSION)

git add VERSION
git commit -m "release version ${version}"
git tag v"${version}"

# NOTE - This push can fail if someone pushed to main while the build
#  was running. We're choosing to mostly ignore this situation due to our
#  currently fairly low commit velocity.
commit_hash=$(git rev-parse --short HEAD)
git push origin ${commit_hash}:main --tags
