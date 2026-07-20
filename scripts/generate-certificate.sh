#!/bin/bash

set -e

project_dir=$(cd "$(dirname "$0")" && pwd)
project_root=$(cd "$project_dir/.." && pwd)
build_dir="$project_root/tests/certs"
certificate_file="$build_dir/certificate.pem"
certificate_key_file="$build_dir/certificate.key"
testdata="$project_root/internal/services/oapi/testdata"

mkdir -p "$build_dir"

needs_regeneration() {
    if [ ! -f "$certificate_file" ] || [ ! -f "$certificate_key_file" ]; then
        return 0
    fi

    if ! openssl x509 -checkend 0 -noout -in "$certificate_file" >/dev/null 2>&1; then
        return 0
    fi

    return 1
}

generate_certificate() {
    openssl req -x509 -sha256 -nodes -newkey rsa:4096 \
        -keyout "$certificate_key_file" \
        -days 1 \
        -out "$certificate_file" \
        -subj /CN=domain.com
}

echo "Ensuring certificates in $build_dir"

if needs_regeneration; then
    rm -f "$certificate_file" "$certificate_key_file"
    generate_certificate
fi

mkdir -p "$testdata"
cp "$certificate_file" "$certificate_key_file" "$testdata/"

exit 0
