#!/usr/bin/env sh
set -eu

REPO_ROOT="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$REPO_ROOT"

go build -o devsnap ./cmd/devenv

for ex in examples/node-npm examples/node-pnpm examples/node-yarn examples/python-req examples/python-poetry; do
  echo "==> $ex"
  ./devsnap init --root "$ex"
  ./devsnap doctor --root "$ex"
  ./devsnap restore --root "$ex"
done

