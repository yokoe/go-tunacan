class Tunacan < Formula
  desc ""
  homepage "https://github.com/yokoe/tunacan"
  url "https://github.com/yokoe/tunacan/releases/download/v0.0.1/tunacan_0.0.1_darwin_amd64.tar.gz"
  version "0.0.1"
  sha256 "4be0704617f7488d80703b4f122e3ddf94a191dc8149a78e3d7f17a6cb44748a"

  def install
    bin.install "tunacan"
  end
end
