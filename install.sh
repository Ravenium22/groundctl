#!/bin/sh
# groundctl installer — https://groundctl.dev
# Usage: curl -fsSL https://get.groundctl.dev | sh
# Optional:
#   GROUNDCTL_VERSION=v1.2.3 curl -fsSL https://get.groundctl.dev | sh
#   curl -fsSL https://get.groundctl.dev | sh -s -- --version v1.2.3
set -e

REPO="groundctl/groundctl"
INSTALL_DIR="/usr/local/bin"
BINARY="ground"

# Colors (when terminal supports it)
RED='\033[0;31m'
GREEN='\033[0;32m'
DIM='\033[0;90m'
BOLD='\033[1m'
NC='\033[0m'

info() { printf "${GREEN}▸${NC} %s\n" "$1"; }
warn() { printf "${RED}▸${NC} %s\n" "$1"; }
dim()  { printf "${DIM}%s${NC}\n" "$1"; }

normalize_version() {
    case "$1" in
        v*) echo "$1" ;;
        *)  echo "v$1" ;;
    esac
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *) warn "Unsupported OS: $(uname -s)"; exit 1 ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)  echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *) warn "Unsupported architecture: $(uname -m)"; exit 1 ;;
    esac
}

# Get latest release tag from GitHub
get_latest_version() {
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        warn "curl or wget is required"
        exit 1
    fi
}

# Download file
download() {
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$1" -o "$2"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$1" -O "$2"
    fi
}

main() {
    REQUESTED_VERSION="${GROUNDCTL_VERSION:-}"

    while [ $# -gt 0 ]; do
        case "$1" in
            --version)
                shift
                if [ -z "${1:-}" ]; then
                    warn "--version requires a value (for example: v1.0.0)"
                    exit 1
                fi
                REQUESTED_VERSION="$1"
                ;;
            -h|--help)
                cat <<EOF
groundctl installer

Usage:
  curl -fsSL https://get.groundctl.dev | sh
  curl -fsSL https://get.groundctl.dev | sh -s -- --version v1.0.0

Environment variable:
  GROUNDCTL_VERSION=v1.0.0
EOF
                exit 0
                ;;
            *)
                warn "Unknown argument: $1"
                exit 1
                ;;
        esac
        shift
    done

    printf "\n${BOLD}  groundctl installer${NC}\n"
    dim "  terraform plan for your local dev machine"
    printf "\n"

    OS=$(detect_os)
    ARCH=$(detect_arch)

    info "Detected: ${OS}/${ARCH}"

    if [ -n "$REQUESTED_VERSION" ]; then
        VERSION=$(normalize_version "$REQUESTED_VERSION")
        info "Requested version: ${VERSION}"
    else
        VERSION=$(get_latest_version)
        if [ -z "$VERSION" ]; then
            warn "Could not determine latest version"
            exit 1
        fi
        info "Latest version: ${VERSION}"
    fi

    # Build download URL
    VERSION_NUM="${VERSION#v}"
    if [ "$OS" = "windows" ]; then
        ARCHIVE="groundctl_${VERSION_NUM}_${OS}_${ARCH}.zip"
    else
        ARCHIVE="groundctl_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
    fi
    URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"

    info "Downloading ${ARCHIVE}..."

    TMPDIR=$(mktemp -d)
    trap "rm -rf $TMPDIR" EXIT

    download "$URL" "$TMPDIR/$ARCHIVE"

    # Extract
    info "Extracting..."
    if [ "$OS" = "windows" ]; then
        unzip -qo "$TMPDIR/$ARCHIVE" -d "$TMPDIR"
    else
        tar -xzf "$TMPDIR/$ARCHIVE" -C "$TMPDIR"
    fi

    # Install
    if [ -w "$INSTALL_DIR" ]; then
        mv "$TMPDIR/$BINARY" "$INSTALL_DIR/$BINARY"
    else
        info "Installing to ${INSTALL_DIR} (requires sudo)..."
        sudo mv "$TMPDIR/$BINARY" "$INSTALL_DIR/$BINARY"
    fi
    chmod +x "$INSTALL_DIR/$BINARY"

    info "Installed groundctl ${VERSION} to ${INSTALL_DIR}/${BINARY}"
    printf "\n"
    dim "  Get started:"
    dim "    ground init       # scan your machine"
    dim "    ground check      # see what's drifted"
    dim "    ground fix        # resolve drift"
    printf "\n"
    dim "  Docs: https://groundctl.dev"
    printf "\n"
}

main "$@"
