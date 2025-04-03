package docpdf

import (
	"math/big"
	"strings"
)

// strHelperGetStringWidth get string width
func strHelperGetStringWidth(str string, fontSize int, ifont iFont) float64 {
	return strHelperGetStringWidthPrecise(str, float64(fontSize), ifont)
}

// strHelperGetStringWidthPrecise get string width with real number fontSize
func strHelperGetStringWidthPrecise(str string, fontSize float64, ifont iFont) float64 {

	w := 0
	bs := []byte(str)
	i := 0
	max := len(bs)
	for i < max {
		w += ifont.GetCw()[bs[i]]
		i++
	}
	return float64(w) * (float64(fontSize) / 1000.0)
}

// createEmbeddedFontSubsetName create Embedded font (subset font) name
func createEmbeddedFontSubsetName(name string) string {
	name = strings.Replace(name, " ", "+", -1)
	name = strings.Replace(name, "/", "+", -1)
	return name
}

// readShortFromByte read short from byte array
func readShortFromByte(data []byte, offset int) (int64, int) {
	buff := data[offset : offset+2]
	num := big.NewInt(0)
	num.SetBytes(buff)
	u := num.Uint64()
	var v int64
	if u >= 0x8000 {
		v = int64(u) - 65536
	} else {
		v = int64(u)
	}
	return v, 2
}

// readUShortFromByte read ushort from byte array
func readUShortFromByte(data []byte, offset int) (uint64, int) {
	buff := data[offset : offset+2]
	num := big.NewInt(0)
	num.SetBytes(buff)
	return num.Uint64(), 2
}
