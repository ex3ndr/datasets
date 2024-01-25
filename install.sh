#!/bin/sh
set -eu

# Helper functions
status() { echo ">>> $*" >&2; }
error() { echo "ERROR $*"; exit 1; }
warning() { echo "WARNING: $*"; }

# Temporary directory
TEMP_DIR=$(mktemp -d)
cleanup() { rm -rf $TEMP_DIR; }
trap cleanup EXIT

# Check if a command is available
available() { command -v $1 >/dev/null; }
require() {
    local MISSING=''
    for TOOL in $*; do
        if ! available $TOOL; then
            MISSING="$MISSING $TOOL"
        fi
    done

    echo $MISSING
}

# Check if the script is running on Linux or MacOS
SYS=$(uname -s)
case "$SYS" in
    Linux) SYS="linux" ;;
    Darwin) SYS="darwin" ;;
    *) error "Unsupported OS: $SYS" ;;
esac

# Check architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) error "Unsupported architecture: $ARCH" ;;
esac

# Resolving sudo command
SUDO=
if [ "$(id -u)" -ne 0 ]; then
    # Running as root, no need for sudo
    if ! available sudo; then
        error "This script requires superuser permissions. Please re-run as root."
    fi

    SUDO="sudo"
fi

# Check if required tools are available
NEEDS=$(require curl)
if [ -n "$NEEDS" ]; then
    status "ERROR: The following tools are required but missing:"
    for NEED in $NEEDS; do
        echo "  - $NEED"
    done
    exit 1
fi

# Download tool
status "Downloading datasets..."
curl --fail --show-error --location --progress-bar -o $TEMP_DIR/datasets "https://github.com/ex3ndr/datasets/releases/latest/download/datasets-$SYS-$ARCH"

# Detecting installation directory
for BINDIR in /usr/local/bin /usr/bin /bin; do
    echo $PATH | grep -q $BINDIR && break || continue
done

# Installing
status "Installing datasets to $BINDIR..."
$SUDO install -o0 -g0 -m755 -d $BINDIR
$SUDO install -o0 -g0 -m755 $TEMP_DIR/datasets $BINDIR/datasets