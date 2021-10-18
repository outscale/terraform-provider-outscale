#!/bin/bash
set -e
apt-get update -y
apt-get upgrade -y
apt-get install -y curl
curl -o /tmp/install.sh "https://install.yunohost.org"
chmod +x /tmp/install.sh
/tmp/install.sh -a
