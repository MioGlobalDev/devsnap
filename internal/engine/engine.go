package engine

import (
	"sort"
)

type Step struct {
	Run string
}

type Detection struct {
	Kind string // "node" | "python"
	Data any
}

type Detector interface {
	Kind() string
	Priority() int
	Detect(root string) (Detection, []Step, bool, error)
}

func DetectAll(root string, detectors []Detector) ([]Detection, []Step, error) {
	sort.SliceStable(detectors, func(i, j int) bool {
		return detectors[i].Priority() < detectors[j].Priority()
	})

	var detections []Detection
	var steps []Step

	for _, d := range detectors {
		det, s, ok, err := d.Detect(root)
		if err != nil {
			return nil, nil, err
		}
		if !ok {
			continue
		}
		detections = append(detections, det)
		steps = append(steps, s...)
	}

	return detections, steps, nil
}

