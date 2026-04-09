# devsnap

Minimal dev environment snapshot tool.

## Usage

go build -o devsnap ./cmd/devenv
./devsnap init
./devsnap doctor
./devsnap restore

## Example (real path)

```bash
git clone <your-project>
cd <your-project>
devsnap init
devsnap doctor
devsnap restore
```

Example output:

```text
[devsnap] Environment check

Summary:
  Node:   ISSUES
  Python: OK

Issues:
  - pnpm not found
    pnpm is not installed; npx fallback is available (already used by devsnap)

Fix:
  - pnpm not found
    - npm install -g pnpm
    - Or keep using: npx -y pnpm@9.15.4 ... (devsnap already falls back)
```

