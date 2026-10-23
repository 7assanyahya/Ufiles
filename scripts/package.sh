#!/usr/bin/env bash
# Builds release archives for each supported platform into dist/.
# Usage: scripts/package.sh <version>   e.g. scripts/package.sh v1.0.0
set -euo pipefail

VERSION="${1:?usage: $0 <version>}"
TARGETS="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64"

rm -rf dist
mkdir -p dist

for target in $TARGETS; do
  os="${target%/*}"
  arch="${target#*/}"
  name="ufiles-${VERSION}-${os}-${arch}"
  dir="dist/${name}"
  bin="ufiles"
  [ "$os" = "windows" ] && bin="ufiles.exe"

  mkdir -p "$dir/static/servfile"
  CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" \
    go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o "$dir/$bin" .
  cp static/*.html static/*.css static/*.js static/*.svg "$dir/static/"
  cp README.md LICENSE "$dir/"

  if [ "$os" = "windows" ]; then
    (cd dist && zip -qr "${name}.zip" "$name")
  else
    tar -C dist -czf "dist/${name}.tar.gz" "$name"
  fi
  rm -rf "$dir"
done

(cd dist && sha256sum ufiles-* > checksums.txt)
ls -lh dist
