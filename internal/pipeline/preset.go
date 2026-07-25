package pipeline

import (
	"strings"

	"github.com/optivor/optivor/internal/config"
)

// ApplyPreset merges preset configuration options into TransformParams.
func ApplyPreset(preset config.PresetConfig, params TransformParams) TransformParams {
	if preset.Width > 0 && params.Width <= 0 {
		params.Width = preset.Width
	}
	if preset.Height > 0 && params.Height <= 0 {
		params.Height = preset.Height
	}
	if preset.Format != "" && params.Format == "" {
		params.Format = preset.Format
	}
	if preset.Fit != "" && params.Fit == "" {
		fitStr := strings.ToLower(preset.Fit)
		switch FitMode(fitStr) {
		case FitCover, FitContain, FitFill, FitSmart:
			params.Fit = FitMode(fitStr)
		}
	}
	return params
}
