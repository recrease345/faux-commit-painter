
set -e

REPO="recrease345/faux-commit-painter"
BIN="gitbrush"

command -v git >/dev/null 2>&1 || { echo "✗ git is required: https://git-scm.com"; exit 1; }


os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)
case "$arch" in
  x86_64|amd64)  arch=amd64 ;;
  aarch64|arm64) arch=arm64 ;;
  *) echo "✗ Unsupported architecture: $arch"; exit 1 ;;
esac
case "$os" in
  linux|darwin) ;;
  *) echo "✗ This script supports Linux/macOS. On Windows download the .exe from the Releases page"; exit 1 ;;
esac


dest="$HOME/.local/bin"
mkdir -p "$dest"
echo "→ downloading $BIN for $os/$arch..."
curl -fsSL "https://github.com/$REPO/releases/latest/download/$BIN-$os-$arch" -o "$dest/$BIN"
chmod +x "$dest/$BIN"


case ":$PATH:" in
  *":$dest:") ;;
  *) echo "⚠  Add to ~/.bashrc (or ~/.zshrc):  export PATH=\"\$HOME/.local/bin:\$PATH\"" ;;
esac

echo "✓ Installed: $dest/$BIN"


if [ "$1" = "run" ]; then
  exec "$dest/$BIN" </dev/tty
fi
