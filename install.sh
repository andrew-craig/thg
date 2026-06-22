#!/usr/bin/env bash
#
# install.sh — build and install the `thg` binary.
#
# Usage:
#   ./install.sh              # build and install to the default location
#   PREFIX=~/.local ./install.sh   # install to ~/.local/bin
#
# Requires: Go 1.25+ and macOS (thg targets Things 3 on macOS).

set -euo pipefail

# Resolve the directory this script lives in, so it works from anywhere.
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Install destination. Override with PREFIX=... (binary goes in $PREFIX/bin).
PREFIX="${PREFIX:-/usr/local}"
BIN_DIR="$PREFIX/bin"

# --- Sanity checks ----------------------------------------------------------

if [[ "$(uname)" != "Darwin" ]]; then
  echo "warning: thg targets Things 3 on macOS; building anyway." >&2
fi

if ! command -v go >/dev/null 2>&1; then
  echo "error: Go is not installed. Install it from https://go.dev/dl/ (1.25+)." >&2
  exit 1
fi

# --- Build ------------------------------------------------------------------

echo "Building thg..."
(cd "$SCRIPT_DIR" && go build -o thg .)

# --- Install ----------------------------------------------------------------

echo "Installing to $BIN_DIR/thg"
if [[ -w "$BIN_DIR" ]] || mkdir -p "$BIN_DIR" 2>/dev/null && [[ -w "$BIN_DIR" ]]; then
  install -m 0755 "$SCRIPT_DIR/thg" "$BIN_DIR/thg"
else
  echo "  (need elevated permissions for $BIN_DIR; using sudo)"
  sudo install -d -m 0755 "$BIN_DIR"
  sudo install -m 0755 "$SCRIPT_DIR/thg" "$BIN_DIR/thg"
fi

echo
echo "Installed: $BIN_DIR/thg"
if ! command -v thg >/dev/null 2>&1; then
  echo "Note: $BIN_DIR is not on your PATH. Add it, e.g.:"
  echo "  export PATH=\"$BIN_DIR:\$PATH\""
fi
echo "Run 'thg --help' to get started."
