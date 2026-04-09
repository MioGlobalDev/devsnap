Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

Push-Location (Split-Path -Parent $MyInvocation.MyCommand.Path)
try {
  $repoRoot = Resolve-Path ".."
  Push-Location $repoRoot
  try {
    go build -o devsnap.exe .\cmd\devenv

    $examples = @(
      "examples\node-npm",
      "examples\node-pnpm",
      "examples\node-yarn",
      "examples\python-req",
      "examples\python-poetry"
    )

    foreach ($ex in $examples) {
      Write-Host "==> $ex"
      .\devsnap.exe init --root $ex
      .\devsnap.exe doctor --root $ex
      .\devsnap.exe restore --root $ex
    }
  } finally {
    Pop-Location
  }
} finally {
  Pop-Location
}

