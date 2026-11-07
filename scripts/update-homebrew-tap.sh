#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 2 ]]; then
  echo "usage: $0 <version> <tap-repo>" >&2
  exit 1
fi

version="${1#v}"
tap_repo="$2"
tap_dir="$(mktemp -d /tmp/homebrew-tagoly.XXXXXX)"

: "${HOMEBREW_TAP_TOKEN:?HOMEBREW_TAP_TOKEN is required}"

artifact_sha() {
  local path="$1"
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$path" | awk '{print $1}'
  else
    shasum -a 256 "$path" | awk '{print $1}'
  fi
}

darwin_arm64_sha="$(artifact_sha dist/Tagoly_Darwin_arm64.tar.gz)"
darwin_x86_64_sha="$(artifact_sha dist/Tagoly_Darwin_x86_64.tar.gz)"
linux_arm64_sha="$(artifact_sha dist/Tagoly_Linux_arm64.tar.gz)"
linux_x86_64_sha="$(artifact_sha dist/Tagoly_Linux_x86_64.tar.gz)"

git clone "https://x-access-token:${HOMEBREW_TAP_TOKEN}@github.com/${tap_repo}.git" "$tap_dir"

cat >"$tap_dir/tagoly.rb" <<EOF
class Tagoly < Formula
  desc "Smart CLI tool to assist with Git commits"
  homepage "https://github.com/meso1007/Tagoly"
  version "$version"

  if OS.mac?
    if Hardware::CPU.intel?
      url "https://github.com/meso1007/Tagoly/releases/download/v$version/Tagoly_Darwin_x86_64.tar.gz"
      sha256 "$darwin_x86_64_sha"
    else
      url "https://github.com/meso1007/Tagoly/releases/download/v$version/Tagoly_Darwin_arm64.tar.gz"
      sha256 "$darwin_arm64_sha"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/meso1007/Tagoly/releases/download/v$version/Tagoly_Linux_arm64.tar.gz"
      sha256 "$linux_arm64_sha"
    else
      url "https://github.com/meso1007/Tagoly/releases/download/v$version/Tagoly_Linux_x86_64.tar.gz"
      sha256 "$linux_x86_64_sha"
    end
  end

  def install
    bin.install "tagoly"
  end

  test do
    system "#{bin}/tagoly", "--version"
  end
end
EOF

cd "$tap_dir"

if git diff --quiet -- tagoly.rb; then
  echo "Homebrew formula is already up to date."
  exit 0
fi

git config user.name "github-actions[bot]"
git config user.email "github-actions[bot]@users.noreply.github.com"
git add tagoly.rb
git commit -m "chore(release): update tagoly formula for v$version"
git push origin HEAD:main
