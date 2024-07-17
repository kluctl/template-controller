#!/bin/sh

set -e

VERSION=$1

VERSION_REGEX='v([0-9]*)\.([0-9]*)\.([0-9]*)'
VERSION_REGEX_SED='v\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\)'

if [ ! -z "$(git status --porcelain)" ]; then
  echo "working directory is dirty!"
  exit 1
fi

if [ -z "$VERSION" ]; then
  echo "No version specified, using 'git sv next-version'"
  VERSION=v$(git sv next-version)
fi

if [[ ! ($VERSION =~ $VERSION_REGEX) ]]; then
  echo "version is invalid"
  exit 1
fi

VERSION_NO_V=$(echo $VERSION | sed 's/$v//g')

echo VERSION=$VERSION
echo VERSION_NO_V=$VERSION_NO_V

cat docs/install.md | sed "s/ref=$VERSION_REGEX_SED/ref=$VERSION/g" > docs/install.md.tmp
mv docs/install.md.tmp docs/install.md

SED_FILES="\
  docs/install.md \
"
ADD_FILES="$SED_FILES"

yq -i ".version=\"$VERSION_NO_V\" | .appVersion=\"$VERSION\"" deploy/charts/template-controller/Chart.yaml
ADD_FILES="$ADD_FILES deploy/charts/template-controller/Chart.yaml"

for f in $SED_FILES; do
  echo "Replacing version in $f"
  cat $f | sed "s/$VERSION_REGEX_SED/$VERSION/g" > $f.tmp
  mv $f.tmp $f
done

for f in $ADD_FILES; do
  git add $f
done

echo "committing"
git commit -o -m "build: Preparing release $VERSION" -- $ADD_FILES

echo "tagging"
git tag -f $VERSION
