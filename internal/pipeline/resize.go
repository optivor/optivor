package pipeline

import (
	"fmt"

	"github.com/davidbyttow/govips/v2/vips"
)

func parseHexColor(hex string) (r, g, b uint8) {
	if hex == "" {
		return 255, 255, 255
	}
	if hex[0] == '#' {
		hex = hex[1:]
	}
	if len(hex) == 6 {
		_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
		if err == nil {
			return r, g, b
		}
	}
	return 255, 255, 255
}

// applyResize performs image resizing based on the specified FitMode.
func applyResize(img *vips.ImageRef, params TransformParams) error {
	targetW := params.Width
	targetH := params.Height

	if targetW <= 0 && targetH <= 0 {
		return nil
	}

	origW := img.Width()
	origH := img.Height()

	if targetW <= 0 {
		targetW = int(float64(origW) * (float64(targetH) / float64(origH)))
	}
	if targetH <= 0 {
		targetH = int(float64(origH) * (float64(targetW) / float64(origW)))
	}

	switch params.Fit {
	case FitCover:
		// Center crop resize
		if err := img.Thumbnail(targetW, targetH, vips.InterestingCentre); err != nil {
			return fmt.Errorf("failed cover thumbnail: %w", err)
		}

	case FitSmart:
		// Smart crop resize using attention-based cropping
		if err := img.Thumbnail(targetW, targetH, vips.InterestingAttention); err != nil {
			return fmt.Errorf("failed smart thumbnail: %w", err)
		}

	case FitFill:
		// Scale directly to target dimensions ignoring aspect ratio
		hScale := float64(targetW) / float64(origW)
		vScale := float64(targetH) / float64(origH)
		if err := img.ResizeWithVScale(hScale, vScale, vips.KernelAuto); err != nil {
			return fmt.Errorf("failed fill resize: %w", err)
		}

	case FitContain:
		// Fit within dimensions preserving aspect ratio, then pad/embed if needed
		if err := img.ThumbnailWithSize(targetW, targetH, vips.InterestingNone, vips.SizeDown); err != nil {
			return fmt.Errorf("failed contain thumbnail: %w", err)
		}

		currentW := img.Width()
		currentH := img.Height()

		if currentW != targetW || currentH != targetH {
			left := (targetW - currentW) / 2
			top := (targetH - currentH) / 2

			if params.Format == "webp" || img.HasAlpha() {
				// Transparent background for formats supporting alpha
				bgColor := &vips.ColorRGBA{R: 0, G: 0, B: 0, A: 0}
				if err := img.EmbedBackgroundRGBA(left, top, targetW, targetH, bgColor); err != nil {
					return fmt.Errorf("failed contain embed rgba: %w", err)
				}
			} else {
				r, g, b := parseHexColor(params.ContainBackgroundColor)
				bgColor := &vips.Color{R: r, G: g, B: b}
				if err := img.EmbedBackground(left, top, targetW, targetH, bgColor); err != nil {
					return fmt.Errorf("failed contain embed: %w", err)
				}
			}
		}

	default:
		// Default to cover if fit mode is unspecified
		if err := img.Thumbnail(targetW, targetH, vips.InterestingCentre); err != nil {
			return fmt.Errorf("failed default cover thumbnail: %w", err)
		}
	}

	return nil
}
