#!/bin/bash

if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Usage: $0 <snapshot_name> <make_target> [args...]"
    exit 1
fi

SNAPSHOT_NAME=$1
MAKE_TARGET=$2
shift 2

RANDOM_SUFFIX=$(tr -dc 'a-z' </dev/urandom | head -c 5)
PROFILE_NAME="${SNAPSHOT_NAME}_${RANDOM_SUFFIX}"

if echo "$MAKE_TARGET" | grep -q "oks"; then
    PROVIDER_TYPE="outscale_oks"
else
    PROVIDER_TYPE="outscale_oapi"
fi

cleanup() {
    local exit_code=$?
    echo "Cleaning up with frieza..."
    frieza clean $SNAPSHOT_NAME --auto-approve 2>/dev/null || true
    frieza snapshot rm $SNAPSHOT_NAME 2>/dev/null || true
    frieza profile rm $PROFILE_NAME 2>/dev/null || true
    exit $exit_code
}

trap cleanup EXIT INT TERM

echo "Creating frieza profile with name: $PROFILE_NAME (provider: $PROVIDER_TYPE)"
if ! frieza profile new $PROVIDER_TYPE $PROFILE_NAME --region=$OSC_REGION --ak=$OSC_ACCESS_KEY --sk=$OSC_SECRET_KEY; then
    echo "ERROR: failed to create frieza profile"
    exit 1
fi

echo "Creating snapshot..."
if ! frieza snapshot new $SNAPSHOT_NAME $PROFILE_NAME; then
    echo "ERROR: failed to create snapshot"
    exit 1
fi

echo "Running tests..."
make $MAKE_TARGET "$@"
