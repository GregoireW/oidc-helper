class OidcHelper < Formula
  desc "OIDC Helper - Lightweight CLI tool to simplify OIDC authentication flows"
  homepage "https://github.com/GregoireW/oidc-helper"
  license "Apache-2.0"

  version "0.1.0"

  on_macos do
    on_arm do
      url "https://github.com/GregoireW/oidc-helper/releases/download/v0.1.0/oidc-helper_0.1.0_darwin_arm64.tar.gz"
      sha256 "PLACEHOLDER"
    end
    on_intel do
      url "https://github.com/GregoireW/oidc-helper/releases/download/v0.1.0/oidc-helper_0.1.0_darwin_amd64.tar.gz"
      sha256 "PLACEHOLDER"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/GregoireW/oidc-helper/releases/download/v0.1.0/oidc-helper_0.1.0_linux_arm64.tar.gz"
      sha256 "PLACEHOLDER"
    end
    on_intel do
      url "https://github.com/GregoireW/oidc-helper/releases/download/v0.1.0/oidc-helper_0.1.0_linux_amd64.tar.gz"
      sha256 "PLACEHOLDER"
    end
  end

  def install
    bin.install "oidc-helper"
  end

  test do
    system "#{bin}/oidc-helper", "--version"
  end
end
