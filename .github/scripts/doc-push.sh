#!/bin/bash

set -e 

if [ "$TAG" == "" ]; then
    echo "We need the tag of the doc"
    exit 1
fi

if [ "$SSH_PRIVATE_KEY" == "" ]; then
    echo "We need the SSH key of the bot"
    exit 1   
fi

echo "$SSH_PRIVATE_KEY" > bot.key
chmod 600 bot.key
GIT_SSH_COMMAND="ssh -i bot.key" git push -f origin "autobuild-Documentation-$TAG"
