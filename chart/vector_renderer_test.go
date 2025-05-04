package chart

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/cdvelop/docpdf/chart/testutil"
	"github.com/cdvelop/docpdf/style"
)

func TestVectorRendererPath(t *testing.T) {
	// replaced new assertions helper

	vr, err := SVG(100, 100)
	testutil.AssertNil(t, err)

	typed, isTyped := vr.(*vectorRenderer)
	testutil.AssertTrue(t, isTyped)

	typed.MoveTo(0, 0)
	typed.LineTo(100, 100)
	typed.LineTo(0, 100)
	typed.Close()
	typed.FillStroke()

	buffer := bytes.NewBuffer([]byte{})
	err = typed.Save(buffer)
	testutil.AssertNil(t, err)

	raw := string(buffer.Bytes())

	testutil.AssertTrue(t, strings.HasPrefix(raw, "<svg"))
	testutil.AssertTrue(t, strings.HasSuffix(raw, "</svg>"))
}

func TestVectorRendererMeasureText(t *testing.T) {
	// replaced new assertions helper

	f, err := GetDefaultFont()
	testutil.AssertNil(t, err)

	vr, err := SVG(100, 100)
	testutil.AssertNil(t, err)

	vr.SetDPI(DefaultDPI)
	vr.SetFont(f)
	vr.SetFontSize(12.0)

	tb := vr.MeasureText("Ljp")
	testutil.AssertEqual(t, 21, tb.Width())
	testutil.AssertEqual(t, 15, tb.Height())
}

func TestCanvasStyleSVG(t *testing.T) {
	// replaced new assertions helper

	f, err := GetDefaultFont()
	testutil.AssertNil(t, err)

	set := Style{
		StrokeColor: style.ColorWhite,
		StrokeWidth: 5.0,
		FillColor:   style.ColorWhite,
		FontColor:   style.ColorWhite,
		Font:        f,
		Padding:     DefaultBackgroundPadding,
	}

	vCanvas := &vectorCanvas{dpi: DefaultDPI}

	svgString := vCanvas.styleAsSVG(set)
	testutil.AssertNotEmpty(t, svgString)
	testutil.AssertTrue(t, strings.HasPrefix(svgString, "style=\""))
	testutil.AssertTrue(t, strings.Contains(svgString, "stroke:rgba(255,255,255,1.0)"))
	testutil.AssertTrue(t, strings.Contains(svgString, "stroke-width:5"))
	testutil.AssertTrue(t, strings.Contains(svgString, "fill:rgba(255,255,255,1.0)"))
	testutil.AssertTrue(t, strings.HasSuffix(svgString, "\""))
}

func TestCanvasClassSVG(t *testing.T) {
	set := Style{
		ClassName: "test-class",
	}

	vCanvas := &vectorCanvas{dpi: DefaultDPI}

	testutil.AssertEqual(t, "class=\"test-class\"", vCanvas.styleAsSVG(set))
}

func TestCanvasCustomInlineStylesheet(t *testing.T) {
	b := strings.Builder{}

	vCanvas := &vectorCanvas{
		w:   &b,
		css: ".background { fill: red }",
	}

	vCanvas.Start(200, 200)

	testutil.AssertContains(t, b.String(), fmt.Sprintf(`<style type="text/css"><![CDATA[%s]]></style>`, vCanvas.css))
}

func TestCanvasCustomInlineStylesheetWithNonce(t *testing.T) {
	b := strings.Builder{}

	vCanvas := &vectorCanvas{
		w:     &b,
		css:   ".background { fill: red }",
		nonce: "RAND0MSTRING",
	}

	vCanvas.Start(200, 200)

	testutil.AssertContains(t, b.String(), fmt.Sprintf(`<style type="text/css" nonce="%s"><![CDATA[%s]]></style>`, vCanvas.nonce, vCanvas.css))
}
