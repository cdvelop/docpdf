package truetype

import (
	"image"
	"image/draw"
	"os"
	"strings"
	"testing"

	"github.com/cdvelop/docpdf/fixedpoint"
	"golang.org/x/image/font"
)

func BenchmarkDrawString(b *testing.B) {
	data, err := os.ReadFile("../licenses/gpl.txt")
	if err != nil {
		b.Fatal(err)
	}
	lines := strings.Split(string(data), "\n")
	data, err = os.ReadFile("../testdata/luxisr.ttf")
	if err != nil {
		b.Fatal(err)
	}
	f, err := Parse(data)
	if err != nil {
		b.Fatal(err)
	}

	dst := image.NewRGBA(image.Rect(0, 0, 800, 600))
	draw.Draw(dst, dst.Bounds(), image.White, image.ZP, draw.Src)
	d := &font.Drawer{
		Dst:  dst,
		Src:  image.Black,
		Face: NewFace(f, nil),
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j, line := range lines {
			d.Dot = fixedpoint.P(0, (j*16)%600)
			d.DrawString(line)
		}
	}
}
