#!/bin/bash

make builder || exit 1
if [ -d "exported" ]; then
	rm -rf "./exported"
fi
mkdir "./exported"
if [ ! -d "radare2-regressions" ]; then
	git clone --depth 2 https://github.com/radare/radare2-regressions || exit 1
fi

find radare2-regressions/new/db/ -type f | while read FNAME; do
	NAME=$(basename "$FNAME")
	./bin/r2r-build "$FNAME" "./exported/$NAME.json" || exit 1
done
