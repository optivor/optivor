package pipeline

import (
	"fmt"

	"github.com/davidbyttow/govips/v2/vips"
)

// exportImage encodes the vips.ImageRef to WebP or native format.
func exportImage(img *vips.ImageRef, format string) ([]byte, string, error) {
	if format == "webp" {
		webpParams := vips.NewWebpExportParams()
		buf, _, err := img.ExportWebp(webpParams)
		if err != nil {
			return nil, "", fmt.Errorf("failed to export webp: %w", err)
		}
		return buf, "image/webp", nil
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
	case vips.ImageTypeGIF:
		contentType = "image/gif"
	case vips.ImageTypeTIFF:
		contentType = "image/tiff"
	case vips.ImageTypeJPEG:
		contentType = "image/jpeg"
	}

	return buf, contentType, nil
}
