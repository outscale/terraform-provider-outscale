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
    PROVIDER_TYPE="oks"
else
    PROVIDER_TYPE="oapi"
fi

cleanup() {
    trap - INT TERM

    echo "Cleaning up with frieza..."
    frieza clean $SNAPSHOT_NAME --auto-approve  || true
    frieza snapshot rm $SNAPSHOT_NAME  || true
    frieza profile rm $PROFILE_NAME  || true
    exit $exit_code
}

echo "Creating frieza profile: $PROFILE_NAME (type: $PROVIDER_TYPE)"
if [ "$PROVIDER_TYPE" = "oapi" ]; then
    if ! frieza profile new outscale_oapi $PROFILE_NAME --region=$OSC_REGION --ak=$OSC_ACCESS_KEY --sk=$OSC_SECRET_KEY; then
        echo "ERROR: failed to create frieza profile"
        exit 1
    fi
    echo "Adding OOS provider to profile..."
    if ! frieza profile add-provider outscale_oos $PROFILE_NAME; then
        echo "ERROR: failed to add OOS provider"
        exit 1
    fi
else
    if ! frieza profile new outscale_oks $PROFILE_NAME --region=$OSC_REGION --ak=$OSC_ACCESS_KEY --sk=$OSC_SECRET_KEY; then
        echo "ERROR: failed to create frieza profile"
        exit 1
    fi
fi

echo "Creating snapshot..."
if ! frieza snapshot new $SNAPSHOT_NAME $PROFILE_NAME; then
    echo "ERROR: failed to create snapshot"
    exit 1
fi

trap cleanup EXIT INT TERM

echo "Running tests..."
make $MAKE_TARGET "$@"
