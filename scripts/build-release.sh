#!/bin/bash

# Build release binaries for grit
# Usage: ./scripts/build-release.sh [VERSION]

set -e

VERSION=${1:-"dev"}
REPO="zpenka/grit"
BUILD_DIR="build/grit-${VERSION}"

echo "Building grit v${VERSION}..."

# Create build directory
mkdir -p ${BUILD_DIR}

# Platforms and architectures
platforms=(
  "linux:amd64"
  "linux:arm64"
  "darwin:amd64"
  "darwin:arm64"
  "windows:amd64"
)

# Build for each platform
for platform in "${platforms[@]}"; do
  IFS=':' read -r os arch <<< "$platform"

  output_name="grit-${version}-${os}-${arch}"
  if [ "$os" = "windows" ]; then
    output_name="${output_name}.exe"
  fi

  echo "Building ${output_name}..."
  GOOS=${os} GOARCH=${arch} go build -o ${BUILD_DIR}/${output_name} .

  # Create SHA256 checksum
  sha256sum ${BUILD_DIR}/${output_name} > ${BUILD_DIR}/${output_name}.sha256
done

echo ""
echo "Build complete! Binaries in ${BUILD_DIR}/"
ls -lh ${BUILD_DIR}/

echo ""
echo "Next steps:"
echo "1. Tag release: git tag v${VERSION}"
echo "2. Create GitHub release with these binaries"
echo "3. Update Homebrew formula with SHA256 and URL"
