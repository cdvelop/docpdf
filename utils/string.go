package utils

import (
	"strconv"
	"strings"
)

// AnyToString converts any value to string without using fmt
// supports: string, int, float64, bool, error, []string, map[string]any types
// for other types, it returns "<nil>" for nil values and "<?>"
func AnyToString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	case error:
		if val != nil {
			return val.Error()
		}
		return "<nil>"
	case []string:
		return "[" + strings.Join(val, ", ") + "]"
	case map[string]any:
		if len(val) == 0 {
			return "{}"
		}
		pairs := make([]string, 0, len(val))
		for k, v := range val {
			pairs = append(pairs, k+": "+AnyToString(v))
		}
		return "{" + strings.Join(pairs, ", ") + "}"
	default:
		// Basic fallback for other types
		if val == nil {
			return "<nil>"
		}
		return "<?>"
	}
}
