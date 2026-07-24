package pipeline

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/optivor/optivor/internal/storage"
	"go.opentelemetry.io/otel"
)

var (
	ErrOversizedImage = errors.New("source image exceeds maximum allowed pixel count")
)

type FitMode string

const (
	FitCover   FitMode = "cover"
	FitContain FitMode = "contain"
	FitFill    FitMode = "fill"
)

type TransformParams struct {
	Width                  int
	Height                 int
	Fit                    FitMode
	Format                 string // "webp" or ""
	ContainBackgroundColor string // e.g. "#ffffff"
	MaxPixels              int
}

var (
	vipsOnce sync.Once
)

// InitVips ensures libvips is initialized for processing.
func InitVips() {
	vipsOnce.Do(func() {
		vips.LoggingSettings(nil, vips.LogLevelError)
		vips.Startup(nil)
	})
}

// ShutdownVips shuts down libvips engine gracefully.
func ShutdownVips() {
	vips.Shutdown()
}

type Pipeline struct{}

func NewPipeline() *Pipeline {
	InitVips()
	return &Pipeline{}
}

// Run executes the fetch -> transform -> encode image pipeline.
// Note: This pipeline is completely unaware of caching (per ADR-0002 & plan.md).
func (p *Pipeline) Run(ctx context.Context, driver storage.StorageDriver, key string, params TransformParams) ([]byte, string, error) {
	ctx, span := otel.Tracer("optivor").Start(ctx, "pipeline.Transform")
	defer span.End()
	// 1. Fetch source object from storage driver
	reader, err := driver.Get(ctx, key)
	if err != nil {
		return nil, "", err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read source image stream: %w", err)
	}

	return p.TransformBytes(ctx, data, params)
}

// TransformBytes transforms raw image bytes directly.
func (p *Pipeline) TransformBytes(ctx context.Context, data []byte, params TransformParams) ([]byte, string, error) {
	ctx, span := otel.Tracer("optivor").Start(ctx, "pipeline.TransformBytes")
	defer span.End()

	// If no transformation or format change is requested, passthrough original image bytes
	if params.Width <= 0 && params.Height <= 0 && params.Format == "" {
		return data, detectContentType(data), nil
	}

	// 2. Transform image using govips
	importParams := vips.NewImportParams()
	img, err := vips.LoadImageFromBuffer(data, importParams)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode source image: %w", err)
	}
	defer img.Close()

	if params.MaxPixels > 0 && (img.Width()*img.Height()) > params.MaxPixels {
		return nil, "", ErrOversizedImage
	}

	if params.Width > 0 || params.Height > 0 {
		if err := applyResize(img, params); err != nil {
			return nil, "", fmt.Errorf("failed to apply resize: %w", err)
		}
	}

	// 3. Encode image
	buf, contentType, err := exportImage(img, params.Format)
	if err != nil {
		return nil, "", fmt.Errorf("failed to encode image: %w", err)
	}

	return buf, contentType, nil
}

func detectContentType(data []byte) string {
	if len(data) >= 8 {
		// PNG signature: \x89PNG\r\n\x1a\n
		if data[0] == 0x89 && data[1] == 'P' && data[2] == 'N' && data[3] == 'G' {
			return "image/png"
		}
		// JPEG signature: \xFF\xD8\xFF
		if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
			return "image/jpeg"
		}
		// GIF signature: GIF87a or GIF89a
		if data[0] == 'G' && data[1] == 'I' && data[2] == 'F' {
			return "image/gif"
		}
		// WebP signature: RIFF....WEBP
		if string(data[0:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
			return "image/webp"
		}
		// AVIF signature: ....ftypavif or ....ftypavis
		if len(data) >= 12 && string(data[4:8]) == "ftyp" && (string(data[8:12]) == "avif" || string(data[8:12]) == "avis") {
			return "image/avif"
		}
	}
	return "application/octet-stream"
}
