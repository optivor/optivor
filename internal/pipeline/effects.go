package pipeline

import (
	"fmt"
	"strings"

	"github.com/davidbyttow/govips/v2/vips"
)

// applyEffects applies dynamic filters (blur, grayscale, pixelate) to the image.
func applyEffects(img *vips.ImageRef, params TransformParams) error {
	if params.Grayscale {
		if err := img.ToColorSpace(vips.InterpretationBW); err != nil {
			return fmt.Errorf("failed to apply grayscale: %w", err)
		}
	}

	if params.Blur > 0 {
		// Cap max blur radius for performance safety
		blurRadius := params.Blur
		if blurRadius > 100 {
			blurRadius = 100
		}
		if err := img.GaussianBlur(blurRadius); err != nil {
			return fmt.Errorf("failed to apply blur: %w", err)
		}
	}

	if params.Pixelate > 1 {
		block := params.Pixelate
		if block > 50 {
			block = 50
		}
		scale := 1.0 / float64(block)
		if err := img.Resize(scale, vips.KernelNearest); err != nil {
			return fmt.Errorf("failed pixelate downscale: %w", err)
		}
		if err := img.Resize(float64(block), vips.KernelNearest); err != nil {
			return fmt.Errorf("failed pixelate upscale: %w", err)
		}
	}

	return nil
}

// applyFocalCrop handles direct coordinate-based focal cropping (focal=x,y).
func applyFocalCrop(img *vips.ImageRef, targetW, targetH int, focalX, focalY float64) error {
	origW := img.Width()
	origH := img.Height()

	if targetW <= 0 {
		targetW = origW
	}
	if targetH <= 0 {
		targetH = origH
	}

	if targetW >= origW && targetH >= origH {
		return nil
	}

	// Calculate target aspect ratio and scale image first so it covers target dimensions
	scaleX := float64(targetW) / float64(origW)
	scaleY := float64(targetH) / float64(origH)
	scale := scaleX
	if scaleY > scaleX {
		scale = scaleY
	}

	if err := img.Resize(scale, vips.KernelAuto); err != nil {
		return fmt.Errorf("failed focal crop resize: %w", err)
	}

	newW := img.Width()
	newH := img.Height()

	// Clamp focal point normalized coordinates between 0 and 1
	if focalX < 0 {
		focalX = 0
	}
	if focalX > 1 {
		focalX = 1
	}
	if focalY < 0 {
		focalY = 0
	}
	if focalY > 1 {
		focalY = 1
	}

	// Center point in pixels after scaling
	centerX := int(float64(newW) * focalX)
	centerY := int(float64(newH) * focalY)

	left := centerX - (targetW / 2)
	top := centerY - (targetH / 2)

	// Clamp bounds
	if left < 0 {
		left = 0
	}
	if top < 0 {
		top = 0
	}
	if left+targetW > newW {
		left = newW - targetW
	}
	if top+targetH > newH {
		top = newH - targetH
	}

	if err := img.ExtractArea(left, top, targetW, targetH); err != nil {
		return fmt.Errorf("failed focal extract area: %w", err)
	}

	return nil
}

// applyOverlay composites an overlay image onto the target image with gravity and opacity.
func applyOverlay(img *vips.ImageRef, params TransformParams) error {
	if len(params.OverlayBytes) == 0 {
		return nil
	}

	overlayImg, err := vips.LoadImageFromBuffer(params.OverlayBytes, nil)
	if err != nil {
		return fmt.Errorf("failed to decode overlay image: %w", err)
	}
	defer overlayImg.Close()

	// Scale overlay relative to base image if requested
	if params.OverlayScale > 0 && params.OverlayScale <= 100 {
		targetOverlayW := int(float64(img.Width()) * (params.OverlayScale / 100.0))
		if targetOverlayW > 0 && targetOverlayW != overlayImg.Width() {
			scale := float64(targetOverlayW) / float64(overlayImg.Width())
			if err := overlayImg.Resize(scale, vips.KernelAuto); err != nil {
				return fmt.Errorf("failed to scale overlay image: %w", err)
			}
		}
	}

	// Calculate position based on gravity
	padding := 10
	baseW, baseH := img.Width(), img.Height()
	overW, overH := overlayImg.Width(), overlayImg.Height()

	left := (baseW - overW) / 2
	top := (baseH - overH) / 2

	gravity := strings.ToLower(params.Gravity)
	switch gravity {
	case "north_west", "top_left":
		left = padding
		top = padding
	case "north", "top":
		left = (baseW - overW) / 2
		top = padding
	case "north_east", "top_right":
		left = baseW - overW - padding
		top = padding
	case "west", "left":
		left = padding
		top = (baseH - overH) / 2
	case "center":
		left = (baseW - overW) / 2
		top = (baseH - overH) / 2
	case "east", "right":
		left = baseW - overW - padding
		top = (baseH - overH) / 2
	case "south_west", "bottom_left":
		left = padding
		top = baseH - overH - padding
	case "south", "bottom":
		left = (baseW - overW) / 2
		top = baseH - overH - padding
	case "south_east", "bottom_right":
		left = baseW - overW - padding
		top = baseH - overH - padding
	}

	if left < 0 {
		left = 0
	}
	if top < 0 {
		top = 0
	}

	// Composite overlay onto main image
	if err := img.Composite(overlayImg, vips.BlendModeOver, left, top); err != nil {
		return fmt.Errorf("failed overlay composite: %w", err)
	}

	return nil
}
