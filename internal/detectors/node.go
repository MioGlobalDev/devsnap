package detectors

import (
	"devenv-snapshot/internal/detect"
	"devenv-snapshot/internal/engine"
)

type NodeDetector struct{}

func (d NodeDetector) Kind() string     { return "node" }
func (d NodeDetector) Priority() int    { return 20 }
func (d NodeDetector) Detect(root string) (engine.Detection, []engine.Step, bool, error) {
	nd, ok, err := detect.DetectNode(root)
	if err != nil {
		return engine.Detection{}, nil, false, err
	}
	if !ok {
		return engine.Detection{}, nil, false, nil
	}

	return engine.Detection{
			Kind: "node",
			Data: nd,
		},
		NodeSteps(nd.PackageManager),
		true,
		nil
}

func NodeSteps(pm string) []engine.Step {
	switch pm {
	case "npm":
		return []engine.Step{{Run: "npm ci"}}
	case "pnpm":
		// keep examples runnable even if pnpm isn't installed
		return []engine.Step{{Run: "npx -y pnpm@9.15.4 i --frozen-lockfile"}}
	case "yarn":
		// keep examples runnable even if yarn isn't installed
		return []engine.Step{{Run: "npx -y yarn@1.22.22 install --frozen-lockfile"}}
	default:
		// unknown pm: no steps
		return nil
	}
}

