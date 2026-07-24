package pipeline

import (
	"testing"

	"github.com/optivor/optivor/internal/config"
)

func TestApplyPreset(t *testing.T) {
	preset := config.PresetConfig{
		Width:  150,
		Height: 150,
		Format: "webp",
		Fit:    "cover",
	}

	initialParams := TransformParams{}
	res := ApplyPreset(preset, initialParams)

	if res.Width != 150 {
		t.Errorf("expected width 150, got %d", res.Width)
	}
	if res.Height != 150 {
		t.Errorf("expected height 150, got %d", res.Height)
	}
	if res.Format != "webp" {
		t.Errorf("expected format webp, got %s", res.Format)
	}
	if res.Fit != FitCover {
		t.Errorf("expected fit cover, got %s", res.Fit)
	}

	// Explicit request params should override preset defaults
	explicitParams := TransformParams{Width: 300}
	res2 := ApplyPreset(preset, explicitParams)
	if res2.Width != 300 {
		t.Errorf("expected width 300 from explicit param, got %d", res2.Width)
	}
}
