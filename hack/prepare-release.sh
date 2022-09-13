#!/bin/sh

set -e

VERSION=$1

VERSION_REGEX='v([0-9]*)\.([0-9]*)\.([0-9]*)'
VERSION_REGEX_SED='v\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\)'

if [ -z "$VERSION" ]; then
  echo "No version specified, using 'git sv next-version'"
  VERSION=v$(git sv next-version)
fi

if [[ ! ($VERSION =~ $VERSION_REGEX) ]]; then
  echo "version is invalid"
  exit 1
fi

echo VERSION=$VERSION

cat docs/install.md | sed "s/ref=$VERSION_REGEX_SED/ref=$VERSION/g" > docs/install.md.tmp
mv docs/install.md.tmp docs/install.md

cat config/manager/kustomization.yaml | sed "s/$VERSION_REGEX_SED/$VERSION/g" > config/manager/kustomization.yaml.tmp
mv config/manager/kustomization.yaml.tmp config/manager/kustomization.yaml

FILES="docs/install.md config/manager/kustomization.yaml"

for f in $FILES; do
  git add $f
done

echo "committing"
git commit -o -m "build: Preparing release $VERSION" -- $FILES

echo "tagging"
git tag -f $VERSION
