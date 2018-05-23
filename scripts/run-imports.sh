#!/bin/bash
SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CURDIR=$(pwd)
bash "$SCRIPTDIR/import-tests.sh"
make
cd "radare2-regressions"
find "$CURDIR/exported" -type f | while read FNAME; do
	"$CURDIR/bin/r2r" "$FNAME"
done
