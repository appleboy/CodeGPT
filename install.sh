#!/usr/bin/env bash

set -euo pipefail

# Create temp directory for downloads.
TMPDIR="$(mktemp -d)"
function cleanup() {
  rm -rf "${TMPDIR}"
}
trap cleanup EXIT INT TERM

APP=CodeGPT

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
ORANGE='\033[38;2;255;140;0m'
NC='\033[0m' # No Color

function print_message() {
  local level=$1
  local message=$2
  local color=""

  case $level in
  info) color="${GREEN}" ;;
  warning) color="${YELLOW}" ;;
  error) color="${RED}" ;;
  esac

  printf "%b\n" "${color}${message}${NC}"
}

function log_error() {
  print_message error "$1" >&2
  exit "$2"
}

function detect_client_info() {
  # Detect WSL or Cygwin as windows
  if [[ "${CLIENT_PLATFORM}" =~ ^(mingw|cygwin|msys)_nt* ]] || grep -qi microsoft /proc/version 2>/dev/null; then
    CLIENT_PLATFORM="windows"
  fi

  case "${CLIENT_PLATFORM}" in
  darwin | linux | windows) ;;
  *) log_error "Unknown or unsupported platform: ${CLIENT_PLATFORM}. Supported platforms are Linux, Darwin, and Windows." 2 ;;
  esac

  case "${CLIENT_ARCH}" in
  x86_64* | i?86_64* | amd64*) CLIENT_ARCH="amd64" ;;
  aarch64* | arm64*) CLIENT_ARCH="arm64" ;;
  armv7l*) CLIENT_ARCH="arm-7" ;;
  *) log_error "Unknown or unsupported architecture: ${CLIENT_ARCH}. Supported architectures are x86_64, i686, arm64, armv7." 3 ;;
  esac
}

function download_and_install() {
  DOWNLOAD_URL_PREFIX="${RELEASE_URL}/v${VERSION}"
  CLIENT_BINARY="CodeGPT-${VERSION}-${CLIENT_PLATFORM}-${CLIENT_ARCH}"
  print_message info "Downloading ${CLIENT_BINARY} from ${DOWNLOAD_URL_PREFIX}"
  mkdir -p "$INSTALL_DIR" || log_error "Failed to create directory: $INSTALL_DIR" 5

  # Use temp dir for download
  TARGET="${TMPDIR}/${CLIENT_BINARY}"

  curl -# -fSL --retry 5 --keepalive-time 2 ${INSECURE_OPTION} "${DOWNLOAD_URL_PREFIX}/${CLIENT_BINARY}" -o "${TARGET}"
  chmod +x "${TARGET}" || log_error "Failed to set executable permission on: ${TARGET}" 7
  # Move the binary to install dir and rename to codegpt
  mv "${TARGET}" "${INSTALL_DIR}/codegpt" || log_error "Failed to move ${TARGET} to ${INSTALL_DIR}/codegpt" 8
  # show the version
  print_message info "Installed ${ORANGE}${CLIENT_BINARY}${NC} to ${GREEN}${INSTALL_DIR}${NC}"
  print_message info "Run ${ORANGE}codegpt version${NC} to show the version"
  print_message info ""
  print_message info "==============================="
  "${INSTALL_DIR}/codegpt" version
  print_message info "==============================="
  print_message info ""
  print_message info "✅ Installation completed successfully!"
}

function add_to_path() {
  local config_file=$1
  local command=$2

  if grep -Fxq "$command" "$config_file"; then
    print_message info "Configuration already exists in $config_file, skipping"
    return 0
  fi

  if [[ -w $config_file ]]; then
    printf "\n# codegpt\n" >>"$config_file"
    echo "$command" >>"$config_file"
    print_message info "Successfully added ${ORANGE}codegpt ${GREEN}to \$PATH in $config_file"
  else
    print_message warning "Manually add the directory to $config_file (or similar):"
    print_message info "  $command"
  fi
}

# Fetch latest release version from GitHub if VERSION is not set
function get_latest_version() {
  local latest
  if command -v jq >/dev/null 2>&1; then
    latest=$(curl $INSECURE_OPTION -# --retry 5 -fSL https://api.github.com/repos/appleboy/CodeGPT/releases/latest | jq -r .tag_name)
  else
    latest=$(curl $INSECURE_OPTION -# --retry 5 -fSL https://api.github.com/repos/appleboy/CodeGPT/releases/latest | grep '"tag_name":' | sed -E 's/.*"tag_name": ?"v?([^"]+)".*/\1/')
  fi
  # Remove leading 'v' if present
  latest="${latest#v}"
  echo "$latest"
}

# Check for required commands
for cmd in curl; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    log_error "Error: $cmd is not installed. Please install $cmd to proceed." 1
  fi
done

CURL_INSECURE="${CURL_INSECURE:-false}"
if [[ "${CURL_INSECURE}" != 'true' && "${CURL_INSECURE}" != 'false' ]]; then
  log_error "CURL_INSECURE must be either 'true' or 'false'" 4
fi
if [[ "${CURL_INSECURE}" == 'true' ]]; then
  print_message warning "WARNING: CURL_INSECURE is set to true. Proceeding with insecure download."
  print_message warning "WARNING: You are bypassing SSL certificate verification. This is insecure and may expose you to man-in-the-middle attacks."
fi
INSECURE_OPTION=""
if [[ "${CURL_INSECURE}" == 'true' ]]; then
  INSECURE_OPTION="--insecure"
fi

if [[ -z "${VERSION:-}" ]]; then
  LATEST_VERSION=$(get_latest_version)
  if [[ -z "$LATEST_VERSION" ]]; then
    log_error "Failed to fetch the latest version from GitHub." 6
  fi
  VERSION="$LATEST_VERSION"
fi

# Check if VERSION is a valid semantic version
if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  log_error "Invalid version format: $VERSION. Expected format: x.y.z" 1
fi

RELEASE_URL="${RELEASE_URL:-https://github.com/appleboy/CodeGPT/releases/download}"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.codegpt/bin}"
CLIENT_PLATFORM="${CLIENT_PLATFORM:-$(uname -s | tr '[:upper:]' '[:lower:]')}"
CLIENT_ARCH="${CLIENT_ARCH:-$(uname -m)}"

detect_client_info
download_and_install

XDG_CONFIG_HOME=${XDG_CONFIG_HOME:-$HOME/.config}

current_shell=$(basename "$SHELL")
case $current_shell in
fish)
  config_files="$HOME/.config/fish/config.fish"
  ;;
zsh)
  config_files="$HOME/.zshrc $HOME/.zshenv $XDG_CONFIG_HOME/zsh/.zshrc $XDG_CONFIG_HOME/zsh/.zshenv"
  ;;
bash)
  config_files="$HOME/.bashrc $HOME/.bash_profile $HOME/.profile $XDG_CONFIG_HOME/bash/.bashrc $XDG_CONFIG_HOME/bash/.bash_profile"
  ;;
ash)
  config_files="$HOME/.ashrc $HOME/.profile /etc/profile"
  ;;
sh)
  config_files="$HOME/.ashrc $HOME/.profile /etc/profile"
  ;;
*)
  # Default case if none of the above matches
  config_files="$HOME/.bashrc $HOME/.bash_profile $XDG_CONFIG_HOME/bash/.bashrc $XDG_CONFIG_HOME/bash/.bash_profile"
  ;;
esac

config_file=""
for file in $config_files; do
  if [[ -f $file ]]; then
    config_file=$file
    break
  fi
done

if [[ -z $config_file ]]; then
  log_error "No config file found for $current_shell. Checked files: ${config_files[@]}" 1
fi

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
  case $current_shell in
  fish)
    add_to_path "$config_file" "fish_add_path $INSTALL_DIR"
    ;;
  zsh)
    add_to_path "$config_file" "export PATH=$INSTALL_DIR:\$PATH"
    ;;
  bash)
    add_to_path "$config_file" "export PATH=$INSTALL_DIR:\$PATH"
    ;;
  ash)
    add_to_path "$config_file" "export PATH=$INSTALL_DIR:\$PATH"
    ;;
  sh)
    add_to_path "$config_file" "export PATH=$INSTALL_DIR:\$PATH"
    ;;
  *)
    print_message warning "Manually add the directory to $config_file (or similar):"
    print_message info "  export PATH=$INSTALL_DIR:\$PATH"
    ;;
  esac
fi

print_message info "To use the command, please restart your terminal or run:"
print_message info "  source $config_file"

if [ -n "${GITHUB_ACTIONS-}" ] && [ "${GITHUB_ACTIONS}" == "true" ]; then
  echo "$INSTALL_DIR" >>$GITHUB_PATH
  print_message info "Added $INSTALL_DIR to \$GITHUB_PATH"
fi
