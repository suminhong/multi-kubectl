class Mk < Formula
  desc "Multi-Kubectl: Run kubectl commands on multiple contexts"
  homepage "https://github.com/suminhong/multi-kubectl"
  url "https://github.com/suminhong/multi-kubectl/archive/refs/tags/v0.0.1.tar.gz"
  sha256 "REPLACE_WITH_SHA256"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w")
  end

  test do
    system "#{bin}/mk", "--help"
  end
end
