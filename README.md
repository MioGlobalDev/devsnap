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
    当前环境未检测到 pnpm，可继续使用 npx fallback（已自动）

Fix:
  - pnpm not found
    - npm install -g pnpm
    - 或继续使用：npx -y pnpm@9.15.4 ...（devsnap 已自动 fallback）
```

