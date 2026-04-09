package detectors

import (
	"devenv-snapshot/internal/detect"
	"devenv-snapshot/internal/engine"
)

type PythonDetector struct{}

func (d PythonDetector) Kind() string  { return "python" }
func (d PythonDetector) Priority() int { return 10 } // run before node

func (d PythonDetector) Detect(root string) (engine.Detection, []engine.Step, bool, error) {
	pd, ok, err := detect.DetectPython(root)
	if err != nil {
		return engine.Detection{}, nil, false, err
	}
	if !ok {
		return engine.Detection{}, nil, false, nil
	}

	var steps []engine.Step
	switch pd.Manager {
	case "poetry":
		// fallback to pip if poetry isn't available yet (like node pnpm/yarn fallback)
		steps = []engine.Step{{Run: "poetry install || (echo [devsnap] poetry not found, fallback to pip && pip install -r requirements.txt)"}}
	case "pip":
		steps = []engine.Step{{Run: "pip install -r requirements.txt"}}
	}

	return engine.Detection{Kind: "python", Data: pd}, steps, true, nil
}

