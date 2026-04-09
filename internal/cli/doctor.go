package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"devenv-snapshot/internal/paths"
	"devenv-snapshot/internal/snapshot"
	"devenv-snapshot/internal/toolchain"

	"github.com/spf13/cobra"
)

type doctorFlags struct {
	root string
}

type issue struct {
	Title string
	Body  string
	Fix   []string
}

func NewDoctorCmd() *cobra.Command {
	var df doctorFlags

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check environment and print actionable fixes",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := filepath.Abs(df.root)
			if err != nil {
				return err
			}

			s, err := snapshot.ReadFile(paths.SnapshotPath(root))
			if err != nil {
				return fmt.Errorf("read snapshot: %w", err)
			}

			var issues []issue

			// Toolchain checks (record + compare only)
			if s.Toolchains.Node != "" {
				if cur, ok, err := toolchain.NodeVersion(cmd.Context()); err == nil && ok && cur != "" && cur != s.Toolchains.Node {
					issues = append(issues, issue{
						Title: "Node version mismatch",
						Body:  fmt.Sprintf("expected: %s\ncurrent:  %s", s.Toolchains.Node, cur),
						Fix: []string{
							fmt.Sprintf("安装/切换 Node 到 %s（如 nvm/fnm；Windows 可用 fnm）", s.Toolchains.Node),
						},
					})
				}
			}
			if s.Toolchains.Python != "" {
				if cur, ok, err := toolchain.PythonVersion(cmd.Context()); err == nil && ok && cur != "" && cur != s.Toolchains.Python {
					issues = append(issues, issue{
						Title: "Python version mismatch",
						Body:  fmt.Sprintf("expected: %s\ncurrent:  %s", s.Toolchains.Python, cur),
						Fix: []string{
							fmt.Sprintf("安装/切换 Python 到 %s（如 pyenv/pyenv-win）", s.Toolchains.Python),
						},
					})
				}
			}

			// Missing command checks (based on detections)
			if s.Detections.Node != nil {
				issues = append(issues, checkCommandIssue("node", "Node", []string{
					"安装 Node（https://nodejs.org/）",
				})...)

				switch strings.ToLower(s.Detections.Node.PackageManager) {
				case "npm":
					issues = append(issues, checkCommandIssue("npm", "npm", []string{
						"npm 通常随 Node 安装；若缺失请重装 Node",
					})...)
				case "pnpm":
					// we have npx fallback, still tell user the clean fix
					if _, err := exec.LookPath("pnpm"); err != nil {
						issues = append(issues, issue{
							Title: "pnpm not found",
							Body:  "当前环境未检测到 pnpm，可继续使用 npx fallback（已自动）",
							Fix: []string{
								"npm install -g pnpm",
								"或继续使用：npx -y pnpm@9.15.4 ...（devsnap 已自动 fallback）",
							},
						})
					}
				case "yarn":
					if _, err := exec.LookPath("yarn"); err != nil {
						issues = append(issues, issue{
							Title: "yarn not found",
							Body:  "当前环境未检测到 yarn，可继续使用 npx fallback（已自动）",
							Fix: []string{
								"npm install -g yarn",
								"或继续使用：npx -y yarn@1.22.22 ...（devsnap 已自动 fallback）",
							},
						})
					}
				}
			}

			if s.Detections.Python != nil {
				issues = append(issues, checkCommandIssue("python", "Python", []string{
					"安装 Python（https://www.python.org/downloads/）",
				})...)

				issues = append(issues, checkCommandIssue("pip", "pip", []string{
					"python -m ensurepip --upgrade",
					"或 python -m pip install -U pip",
				})...)

				switch strings.ToLower(s.Detections.Python.Manager) {
				case "poetry":
					if _, err := exec.LookPath("poetry"); err != nil {
						issues = append(issues, issue{
							Title: "poetry not found",
							Body:  "将会 fallback 到 pip（若快照 step 启用了 fallback）",
							Fix: []string{
								"python -m pip install --user poetry",
								"参考：https://python-poetry.org/docs/",
							},
						})
					}
				}
			}

			// Output
			fmt.Fprintln(os.Stdout, "[devsnap] Environment check")
			fmt.Fprintln(os.Stdout, "")

			fmt.Fprintln(os.Stdout, "Summary:")
			fmt.Fprintf(os.Stdout, "  Node:   %s\n", statusLine(s.Detections.Node != nil, hasIssuesKind(issues, "node")))
			fmt.Fprintf(os.Stdout, "  Python: %s\n", statusLine(s.Detections.Python != nil, hasIssuesKind(issues, "python")))
			fmt.Fprintln(os.Stdout, "")

			if len(issues) == 0 {
				fmt.Fprintln(os.Stdout, "Issues: none")
				return nil
			}

			fmt.Fprintln(os.Stdout, "Issues:")
			for _, it := range issues {
				fmt.Fprintf(os.Stdout, "  - %s\n", it.Title)
				if strings.TrimSpace(it.Body) != "" {
					for _, line := range strings.Split(it.Body, "\n") {
						fmt.Fprintf(os.Stdout, "    %s\n", line)
					}
				}
			}
			fmt.Fprintln(os.Stdout, "")

			fmt.Fprintln(os.Stdout, "Fix:")
			for _, it := range issues {
				if len(it.Fix) == 0 {
					continue
				}
				fmt.Fprintf(os.Stdout, "  - %s\n", it.Title)
				for _, f := range it.Fix {
					fmt.Fprintf(os.Stdout, "    - %s\n", f)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&df.root, "root", ".", "project root directory")
	return cmd
}

func checkCommandIssue(bin string, title string, fix []string) []issue {
	if _, err := exec.LookPath(bin); err == nil {
		return nil
	}
	return []issue{{
		Title: fmt.Sprintf("%s not found", title),
		Body:  fmt.Sprintf("missing command: %s", bin),
		Fix:   fix,
	}}
}

func hasIssuesKind(items []issue, kind string) bool {
	for _, it := range items {
		l := strings.ToLower(it.Title)
		switch kind {
		case "node":
			if strings.Contains(l, "node") || strings.Contains(l, "npm") || strings.Contains(l, "pnpm") || strings.Contains(l, "yarn") {
				return true
			}
		case "python":
			if strings.Contains(l, "python") || strings.Contains(l, "pip") || strings.Contains(l, "poetry") {
				return true
			}
		}
		if strings.HasPrefix(l, kind+" ") || strings.HasPrefix(l, kind+":") {
			return true
		}
	}
	return false
}

func statusLine(detected bool, hasIssue bool) string {
	if !detected {
		return "SKIP"
	}
	if hasIssue {
		return "ISSUES"
	}
	return "OK"
}

