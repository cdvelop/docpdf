package docpdf

import (
	"strconv"
	//"fmt"
	"bytes"
)

// fontConvertHelperCw2Str converts main ASCII characters of a FontCW to a string.
func fontConvertHelperCw2Str(cw fontCw) string {
	buff := new(bytes.Buffer)
	buff.WriteString(" ")
	i := 32
	for i <= 255 {
		buff.WriteString(strconv.Itoa(cw[byte(i)]) + " ")
		i++
	}
	return buff.String()
}

// fontConvertHelper_Cw2Str converts main ASCII characters of a FontCW to a string. (for backward compatibility)
// Deprecated: Use fontConvertHelperCw2Str(cw fontCw) instead
func fontConvertHelper_Cw2Str(cw fontCw) string {
	return fontConvertHelperCw2Str(cw)
}
