#!/bin/bash

USAGE="Usage: build.sh -o|--output NAME [--basedir terraform-provider-onefuse] [--darwin] [--linux] [--windows] [--sha256sum]"
GOARCH="amd64"
BASEDIR="terraform-provider-onefuse"

while test $# -gt 0; do
    case "$1" in
        -o|--output)
            shift;
            OUTPUT=$1;
            shift;
            ;;
        -b|--basedir)
            shift;
            BASEDIR=$1;
            shift;
            ;;
        --darwin)
            DARWIN=darwin;
            shift;
            ;;
        --linux)
            LINUX=linux;
            shift;
            ;;
        --windows)
            WINDOWS=windows;
            shift;
            ;;
        --sha256sum)
            SHA256SUM=1;
            shift;
            ;;
        *)
            echo $USAGE;
            exit 1;
            ;;
    esac
done

# Output name or OS was not provided, exit with usage details
if [[ -z $OUTPUT || (-z $DARWIN && -z $LINUX && -z $WINDOWS) ]];
then
    echo $USAGE;
    exit 1;
fi

build_os_pkg() {
    OS=$1;

    echo "Building $OS";

    OUTPUT_PATH="$BASEDIR/$OS/$OUTPUT"
    SHA256_PATH="$BASEDIR/$OS/$OUTPUT.sha256sum"

    mkdir -p "$BASEDIR/$OS"
    GOOS=$OS go build -o $OUTPUT_PATH

    if [ $SHA256SUM ];
    then
        echo "Generating SHA256 Checksum for $OS"
        sha256sum $OUTPUT_PATH > $SHA256_PATH
    fi
}

if [ $DARWIN ];
then
    build_os_pkg $DARWIN;
fi

if [ $LINUX ];
then
    build_os_pkg $LINUX;
fi

if [ $WINDOWS ];
then
    build_os_pkg $WINDOWS;
fi

