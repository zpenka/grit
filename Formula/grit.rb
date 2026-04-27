class Grit < Formula
  desc "Terminal UI for exploring git history"
  homepage "https://github.com/zpenka/grit"
  url "https://github.com/zpenka/grit/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "0000000000000000000000000000000000000000000000000000000000000000"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags",
           "-s -w -X main.Version=#{version}",
           "-o", "grit", "."
    bin.install "grit"
  end

  test do
    system "#{bin}/grit", "--version"
  end
end
