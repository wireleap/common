#!/bin/sh
set -e

SCRIPT_NAME="$(basename "$0")"

fatal() { echo "FATAL [$SCRIPT_NAME]: $*" 1>&2; exit 1; }
info() { echo "INFO [$SCRIPT_NAME]: $*"; }

command -v go >/dev/null || fatal "go not installed"

SRCDIR="$(dirname "$(dirname "$(realpath "$0")")")"
GITVERSION="$($SRCDIR/contrib/gitversion.sh)"

NPROC=
if command -v nproc >/dev/null; then
    NPROC="$( nproc )"
elif command -v grep >/dev/null; then
    NPROC="$( grep -c processor /proc/cpuinfo )"
fi

if [ "$NPROC" -lt 2 ]; then
    NPROC=2
fi

info "running at most $NPROC tests in parallel"

info "testing ..."
go test \
    -parallel "$NPROC" \
    -ldflags "
        $VERSIONS
    " \
    "$@" ./...
