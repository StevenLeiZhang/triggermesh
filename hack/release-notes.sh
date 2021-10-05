#!/usr/bin/env bash

RELEASE=${1:-${GIT_TAG}}
RELEASE=${RELEASE:-${CIRCLE_TAG}}

if [ -z "${RELEASE}" ]; then
  echo "Usage:"
  echo "release-notes.sh VERSION"
  exit 1
fi

if ! git rev-list ${RELEASE} >/dev/null 2>&1; then
	echo "${RELEASE} does not exist"
	exit
fi

KREPO="triggermesh"
BASE_URL="https://github.com/triggermesh/${KREPO}/releases/download/${RELEASE}"
PREV_RELEASE=${PREV_RELEASE:-$(git describe --tags --abbrev=0 ${RELEASE}^ 2>/dev/null)}
PREV_RELEASE=${PREV_RELEASE:-$(git rev-list --max-parents=0 ${RELEASE}^ 2>/dev/null)}
CHANGELOG=$(git log --no-merges --pretty=format:'- [%h] %s (%aN)' ${PREV_RELEASE}..${RELEASE})
if [ $? -ne 0 ]; then
  echo "Error creating changelog"
  exit 1
fi

IMAGE_REPO=$(sed -n -e 's/^IMAGE_REPO[[:space:]].*=[[:space:]]*\(.*\)$/\1/p' Makefile)
COMMANDS=$(ls cmd)
PLATFORMS=$(sed -n -e "s/^\(TARGETS[[:space:]]*?=[[:space:]]*\)\(.*\)$/\2/p" Makefile)
RELEASE_ASSETS_TABLE=$(
  echo -n "| component | artifacts |"; echo
  echo -n "| -- | -- |"; echo
  for command in ${COMMANDS}; do
    echo -n "| ${command} | ([container](https://${IMAGE_REPO}/${command}:${RELEASE}))"
    for platform in ${PLATFORMS}; do
      echo -n " ([${platform}](${BASE_URL}/${command}-${platform%/*}-${platform#*/}))"
    done
    echo -n " |"; echo
  done
  echo
)

cat <<EOF
## Installation

Download TriggerMesh Open Source Components ${RELEASE}

${RELEASE_ASSETS_TABLE}

## Changelog

${CHANGELOG}
EOF