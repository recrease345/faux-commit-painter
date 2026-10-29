
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
  *)
    is_fish=0; is_zsh=0
    [ -n "$FISH_VERSION" ] && is_fish=1
    case "$SHELL" in *fish*) is_fish=1;; *zsh*) is_zsh=1;; esac
    if [ "$is_fish" -eq 1 ]; then
      echo "⚠  fish detected. Run once:   fish_add_path $HOME/.local/bin"
    elif [ "$is_zsh" -eq 1 ]; then
      echo "⚠  Add to ~/.zshrc:   export PATH=\"\$HOME/.local/bin:\$PATH\""
    else
      echo "⚠  Add to ~/.bashrc:  export PATH=\"\$HOME/.local/bin:\$PATH\""
    fi
    ;;
esac

echo "✓ Installed: $dest/$BIN"


if [ "$1" = "run" ]; then
  exec "$dest/$BIN" </dev/tty
fi