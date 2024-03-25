#!/usr/bin/env sh
set -x

export PATH=$PATH:/terrad/build/terrad
BINARY=/terrad/build/terrad
ID=${ID:-0}
LOG=${LOG:-terrad.log}

if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found."
	exit 1
fi

export TERRAD_HOME="/terrad/data/node${ID}/simd"

if [ -d "$(dirname "${TERRAD_HOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${TERRAD_HOME}" "$@" | tee "${TERRAD_HOME}/${LOG}"
else
  "${BINARY}" --home "${TERRAD_HOME}" "$@"
fi
