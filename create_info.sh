#!/bin/bash

VERSION="$1"
BUILD="$2"
RELEASE_DATE="$3"
TERRAFORM_PROVIDER_BIN_FILE_PATH="$4"
CHECKSUM_LINUX="$5"
CHECKSUM_DARWIN="$6"
CHECKSUM_WINDOWS="$7"
cat <<EOF
{
    "version": "$VERSION",
    "build": "$BUILD",
    "linux": "http://downloads.cloudbolt.io/OneFuse/Terraform/$VERSION/linux/$TERRAFORM_PROVIDER_BIN_FILE_PATH",
    "darwin": "http://downloads.cloudbolt.io/OneFuse/Terraform/$VERSION/darwin/$TERRAFORM_PROVIDER_BIN_FILE_PATH",
    "windows": "http://downloads.cloudbolt.io/OneFuse/Terraform/$VERSION/windows/$TERRAFORM_PROVIDER_BIN_FILE_PATH.exe",
    "release_date": "$RELEASE_DATE",
    "checksum_linux": "$CHECKSUM_LINUX",
    "checksum_darwin": "$CHECKSUM_DARWIN",
    "checksum_windows": "$CHECKSUM_WINDOWS"
}
EOF
