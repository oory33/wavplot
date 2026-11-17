package plt

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image/png"
	"os"
	"testing"

	"github.com/youpy/go-wav"
)

// Exercise the font and raster backends used by every public plotting function.
// Keep generated images out of the repository's output directory.
func TestPNGOutput(t *testing.T) {
	fixture, err := os.ReadFile("input/H4.wav")
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(t.TempDir())
	if err := os.Mkdir("output", 0755); err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name string
		draw func()
	}{
		{"wav", func() { Wav(wav.NewReader(bytes.NewReader(fixture)), 0, "wav") }},
		{"spec", func() { Spec(wav.NewReader(bytes.NewReader(fixture)), 0, "spec") }},
		{"complex", func() { Complex([]complex128{1, 2 + 3i, 4 - 2i, 3, 1i}, "complex") }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.draw()
			data, err := os.ReadFile(fmt.Sprintf("output/%s.png", tc.name))
			if err != nil {
				t.Fatal(err)
			}
			img, err := png.Decode(bytes.NewReader(data))
			if err != nil {
				t.Fatal(err)
			}
			if got := img.Bounds(); got.Dx() != 960 || got.Dy() != 960 {
				t.Fatalf("unexpected image dimensions: %v", got)
			}
			first := img.At(0, 0)
			var nonuniform bool
			for y := 0; y < 960 && !nonuniform; y++ {
				for x := 0; x < 960; x++ {
					if img.At(x, y) != first {
						nonuniform = true
						break
					}
				}
			}
			if !nonuniform {
				t.Fatal("plot is blank")
			}
			t.Logf("PNG SHA256: %x", sha256.Sum256(data))
		})
	}
}
