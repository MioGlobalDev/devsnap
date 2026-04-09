Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$dist = Join-Path $repoRoot "dist"
New-Item -ItemType Directory -Force -Path $dist | Out-Null

Push-Location $repoRoot
try {
  $out = Join-Path $dist "devsnap-windows-amd64.exe"
  go build -trimpath -ldflags "-s -w" -o $out .\cmd\devenv

  $hash = (Get-FileHash -Algorithm SHA256 $out).Hash.ToLower()
  $hashLine = "$hash  $(Split-Path -Leaf $out)"
  Set-Content -Path (Join-Path $dist "SHA256SUMS.txt") -Value $hashLine -NoNewline

  Write-Host "built: $out"
  Write-Host "sha256: $hash"
} finally {
  Pop-Location
}

