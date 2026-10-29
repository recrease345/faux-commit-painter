set -e

REPO="recrease345/gitbrush" 
BIN="gitbrush"

command -v git >/dev/null || { echo "Git required: https://git-scm.com"; exit 1; }

# определяем платформу
os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)
case "$arch" in
  x86_64|amd64)   arch=amd64 ;;
  aarch64|arm64)  arch=arm64 ;;
  *) echo "unknown architecture: $arch"; exit 1 ;;
esac
case "$os" in
  linux|darwin) ;;
  *) echo "A script only for linux/macOS, if you are using Windows(install linux pls) download the .exe file from the Releases"; exit 1 ;;
esac

dest="$HOME/.local/bin"
mkdir -p "$dest"
curl -fsSL "https://github.com/$REPO/releases/latest/download/$BIN-$os-$arch" -o "$dest/$BIN"
chmod +x "$dest/$BIN"

case ":$PATH:" in
  *":$dest:") ;;
  *) echo "⚠  add to ~/.bashrc (or ~/.zshrc):  export PATH=\"\$HOME/.local/bin:\$PATH\"" ;;
esac

echo "✓ Installed: $dest/$BIN"

if [ "$1" = "run" ]; then
  exec "$dest/$BIN" </dev/tty
fi