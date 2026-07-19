package pipeline_test

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"testing"

	"github.com/davidbyttow/govips/v2/vips"

	"github.com/optivor/optivor/internal/pipeline"
	"github.com/optivor/optivor/internal/storage"
)

type mockStorageDriver struct {
	data map[string][]byte
}

func (m *mockStorageDriver) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	data, ok := m.data[key]
	if !ok {
		return nil, storage.ErrNotFound
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func createTestJPEG(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, color.RGBA{R: 200, G: 100, B: 50, A: 255})
		}
	}
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	return buf.Bytes()
}

func TestPipeline_Passthrough(t *testing.T) {
	rawJPG := createTestJPEG(100, 100)
	driver := &mockStorageDriver{
		data: map[string][]byte{
			"sample.jpg": rawJPG,
		},
	}

	pipe := pipeline.NewPipeline()
	ctx := context.Background()

	res, contentType, err := pipe.Run(ctx, driver, "sample.jpg", pipeline.TransformParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if contentType != "image/jpeg" {
		t.Errorf("expected image/jpeg, got %s", contentType)
	}
	if !bytes.Equal(res, rawJPG) {
		t.Error("expected raw bytes passthrough to match original")
	}
}

func TestPipeline_ResizeWebP(t *testing.T) {
	rawJPG := createTestJPEG(400, 300)
	driver := &mockStorageDriver{
		data: map[string][]byte{
			"sample.jpg": rawJPG,
		},
	}

	pipe := pipeline.NewPipeline()
	ctx := context.Background()

	modes := []pipeline.FitMode{pipeline.FitCover, pipeline.FitContain, pipeline.FitFill}
	for _, mode := range modes {
		t.Run(string(mode), func(t *testing.T) {
			params := pipeline.TransformParams{
				Width:                  150,
				Height:                 150,
				Fit:                    mode,
				Format:                 "webp",
				ContainBackgroundColor: "#00ff00",
			}

			res, contentType, err := pipe.Run(ctx, driver, "sample.jpg", params)
			if err != nil {
				t.Fatalf("fit mode %s failed: %v", mode, err)
			}

			if contentType != "image/webp" {
				t.Errorf("expected image/webp, got %s", contentType)
			}

			// Verify dimensions using govips
			importParams := vips.NewImportParams()
			vipsImg, err := vips.LoadImageFromBuffer(res, importParams)
			if err != nil {
				t.Fatalf("failed to decode webp output: %v", err)
			}
			defer vipsImg.Close()

			if vipsImg.Width() != 150 || vipsImg.Height() != 150 {
				t.Errorf("expected output 150x150, got %dx%d", vipsImg.Width(), vipsImg.Height())
			}
		})
	}
}

func TestPipeline_OversizedImage(t *testing.T) {
	rawJPG := createTestJPEG(200, 200) // 40,000 pixels
	driver := &mockStorageDriver{
		data: map[string][]byte{
			"sample.jpg": rawJPG,
		},
	}

	pipe := pipeline.NewPipeline()
	ctx := context.Background()

	params := pipeline.TransformParams{
		Width:     100,
		Height:    100,
		MaxPixels: 10000, // lower than 40,000
	}

	_, _, err := pipe.Run(ctx, driver, "sample.jpg", params)
	if err == nil {
		t.Fatal("expected error for oversized image, got nil")
	}
	if err != pipeline.ErrOversizedImage {
		t.Errorf("expected ErrOversizedImage, got %v", err)
	}
}

func BenchmarkPipeline_ResizeWebP(b *testing.B) {
	rawJPG := createTestJPEG(800, 600)
	driver := &mockStorageDriver{
		data: map[string][]byte{
			"sample.jpg": rawJPG,
		},
	}

	pipe := pipeline.NewPipeline()
	ctx := context.Background()
	params := pipeline.TransformParams{
		Width:  200,
		Height: 200,
		Fit:    pipeline.FitCover,
		Format: "webp",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = pipe.Run(ctx, driver, "sample.jpg", params)
	}
}
