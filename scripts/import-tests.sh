#!/bin/bash

make builder || exit 1
if [ -d "exported" ]; then
	rm -rf "./exported"
fi
mkdir "./exported"
if [ ! -d "radare2-regressions" ]; then
	git clone https://github.com/radare/radare2-regressions || exit 1
fi

find radare2-regressions/new/db/ -type f | while read FNAME; do
	NAME=$(basename "$FNAME")
	./bin/r2r-build "$FNAME" "./exported/$NAME.json" || exit 1
done

DIR=$(pwd)
make
cd "radare2-regressions"
find "$DIR/exported" -type f | while read FNAME; do
	"$DIR/bin/r2r" "$FNAME"
done
