package docpdf

import "github.com/cdvelop/docpdf/errs"

// convertNumericToFloat64 : accept numeric types, return float64-value
func convertNumericToFloat64(size any) (fontSize float64, err error) {
	switch size := size.(type) {
	case float32:
		return float64(size), nil
	case float64:
		return float64(size), nil
	case int:
		return float64(size), nil
	case int16:
		return float64(size), nil
	case int32:
		return float64(size), nil
	case int64:
		return float64(size), nil
	case int8:
		return float64(size), nil
	case uint:
		return float64(size), nil
	case uint16:
		return float64(size), nil
	case uint32:
		return float64(size), nil
	case uint64:
		return float64(size), nil
	case uint8:
		return float64(size), nil
	default:
		return 0.0, errs.New("fontSize must be of type (u)int* or float*, not", size)
	}
}

// Determine if a value is numeric
// Courtesy of https://github.com/syyongx/php2go/blob/master/php.go
func is_numeric(val any) bool {
	switch v := val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
	case float32, float64, complex64, complex128:
		return true
	case string:
		str := v
		if str == "" {
			return false
		}
		// Trim any whitespace
		str = customTrimSpace(str)
		//fmt.Println(str)
		if str[0] == '-' || str[0] == '+' {
			if len(str) == 1 {
				return false
			}
			str = str[1:]
		}
		// hex
		if len(str) > 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X') {
			for _, h := range str[2:] {
				if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'f') || (h >= 'A' && h <= 'F')) {
					return false
				}
			}
			return true
		}
		// 0-9,point,Scientific
		p, s, l := 0, 0, len(str)
		for i, v := range str {
			if v == '.' { // point
				if p > 0 || s > 0 || i+1 == l {
					return false
				}
				p = i
			} else if v == 'e' || v == 'E' { // Scientific
				if i == 0 || s > 0 || i+1 == l {
					return false
				}
				s = i
			} else if v < '0' || v > '9' {
				return false
			}
		}
		return true
	}

	return false
}

func in_array(needle any, hystack any) bool {
	switch key := needle.(type) {
	case string:
		for _, item := range hystack.([]string) {
			if key == item {
				return true
			}
		}
	case int:
		for _, item := range hystack.([]int) {
			if key == item {
				return true
			}
		}
	case int64:
		for _, item := range hystack.([]int64) {
			if key == item {
				return true
			}
		}
	default:
		return false
	}
	return false
}

// Taken from png library

// intSize is either 32 or 64.
const intSize = 32 << (^uint(0) >> 63)

func abs(x int) int {
	// m := -1 if x < 0. m := 0 otherwise.
	m := x >> (intSize - 1)

	// In two's complement representation, the negative number
	// of any number (except the smallest one) can be computed
	// by flipping all the bits and add 1. This is faster than
	// code with a branch.
	// See Hacker's Delight, section 2-4.
	return (x ^ m) - m
}

// filterPaeth applies the Paeth filter to the cdat slice.
// cdat is the current row's data, pdat is the previous row's data.
func filterPaeth(cdat, pdat []byte, bytesPerPixel int) {
	var a, b, c, pa, pb, pc int
	for i := 0; i < bytesPerPixel; i++ {
		a, c = 0, 0
		for j := i; j < len(cdat); j += bytesPerPixel {
			b = int(pdat[j])
			pa = b - c
			pb = a - c
			pc = abs(pa + pb)
			pa = abs(pa)
			pb = abs(pb)
			if pa <= pb && pa <= pc {
				// No-op.
			} else if pb <= pc {
				a = b
			} else {
				a = c
			}
			a += int(cdat[j])
			a &= 0xff
			cdat[j] = uint8(a)
			c = b
		}
	}
}

// customTrimSpace quita espacios en blanco al inicio y al final
func customTrimSpace(s string) string {
	i, j := 0, len(s)-1
	// Avanza desde la izquierda
	for i <= j && (s[i] == ' ' || s[i] == '\t' || s[i] == '\r' || s[i] == '\n') {
		i++
	}
	// Avanza desde la derecha
	for j >= i && (s[j] == ' ' || s[j] == '\t' || s[j] == '\r' || s[j] == '\n') {
		j--
	}
	return s[i : j+1]
}

// funcKernOverride  return your custome pair value
type funcKernOverride func(
	leftRune rune,
	rightRune rune,
	leftPair uint,
	rightPair uint,
	pairVal int16,
) int16
