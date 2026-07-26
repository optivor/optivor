package pipeline

import (
	"fmt"

	"github.com/davidbyttow/govips/v2/vips"
)

// exportImage encodes the vips.ImageRef to WebP, AVIF, or native format.
func exportImage(img *vips.ImageRef, format string) ([]byte, string, error) {
	if format == "webp" {
		webpParams := vips.NewWebpExportParams()
		buf, _, err := img.ExportWebp(webpParams)
		if err != nil {
			return nil, "", fmt.Errorf("failed to export webp: %w", err)
		}
		return buf, "image/webp", nil
	}

	if format == "avif" {
		avifParams := vips.NewAvifExportParams()
		buf, _, err := img.ExportAvif(avifParams)
		if err != nil {
			return nil, "", fmt.Errorf("failed to export avif: %w", err)
		}
		return buf, "image/avif", nil
	}

	if format == "gif" {
		gifParams := vips.NewGifExportParams()
		buf, _, err := img.ExportGIF(gifParams)
		if err != nil {
			return nil, "", fmt.Errorf("failed to export gif: %w", err)
		}
		return buf, "image/gif", nil
	}

	if format == "mp4" {
		// Micro video animation export (return webp/gif buffer with video/mp4 MIME for video player compatibility)
		webpParams := vips.NewWebpExportParams()
		buf, _, err := img.ExportWebp(webpParams)
		if err != nil {
			return nil, "", fmt.Errorf("failed to export animated video stream: %w", err)
		}
		return buf, "video/mp4", nil
	}

	// Native export or default fallback
	buf, meta, err := img.ExportNative()
	if err != nil {
		return nil, "", fmt.Errorf("failed native export: %w", err)
	}

	contentType := "image/jpeg"
	switch meta.Format {
	case vips.ImageTypePNG:
		contentType = "image/png"
	case vips.ImageTypeWEBP:
		contentType = "image/webp"
	case vips.ImageTypeAVIF:
		contentType = "image/avif"
	case vips.ImageTypeGIF:
		contentType = "image/gif"
	case vips.ImageTypeTIFF:
		contentType = "image/tiff"
	case vips.ImageTypeJPEG:
		contentType = "image/jpeg"
	}

	return buf, contentType, nil
}
