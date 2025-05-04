package fontengine

import "github.com/cdvelop/docpdf/errs"

var eRROR_NO_KEY_FOUND = errs.New("no key found")
var eRROR_NO_GET_WRONG_TYPE = errs.New("get wrong type")

type ttfInfo map[string]any

func (t ttfInfo) PushString(key string, val string) {
	t[key] = val
}

func (t ttfInfo) PushBytes(key string, val []byte) {
	t[key] = val
}

func (t ttfInfo) PushInt64(key string, val int64) {
	t[key] = val
}

func (t ttfInfo) PushInt(key string, val int) {
	t[key] = val
}

func (t ttfInfo) PushUInt64(key string, val uint) {
	t[key] = val
}

func (t ttfInfo) PushBool(key string, val bool) {
	t[key] = val
}

func (t ttfInfo) PushInt64s(key string, val []int) {
	t[key] = val
}

func (t ttfInfo) PushMapIntInt64(key string, val map[int]int) {
	t[key] = val
}

func (t ttfInfo) GetBool(key string) (bool, error) {
	if val, ok := t[key]; ok {

		if m, ok := val.(bool); ok {
			/* act on str */
			return m, nil
		} else {
			return false, eRROR_NO_GET_WRONG_TYPE
		}
	} else {
		return false, eRROR_NO_KEY_FOUND
	}
}

func (t ttfInfo) GetString(key string) (string, error) {
	if val, ok := t[key]; ok {

		if m, ok := val.(string); ok {
			/* act on str */
			return m, nil
		} else {
			return "", eRROR_NO_GET_WRONG_TYPE
		}
	} else {
		return "", eRROR_NO_KEY_FOUND
	}
}

func (t ttfInfo) GetInt64(key string) (int, error) {
	if val, ok := t[key]; ok {

		if m, ok := val.(int); ok {
			/* act on str */
			return m, nil
		} else {
			return 0, eRROR_NO_GET_WRONG_TYPE
		}
	} else {
		return 0, eRROR_NO_KEY_FOUND
	}
}

func (t ttfInfo) GetInt64s(key string) ([]int, error) {
	if val, ok := t[key]; ok {

		if m, ok := val.([]int); ok {
			/* act on str */
			return m, nil
		} else {
			return nil, eRROR_NO_GET_WRONG_TYPE
		}
	} else {
		return nil, eRROR_NO_KEY_FOUND
	}
}

func (t ttfInfo) GetMapIntInt64(key string) (map[int]int, error) {
	if val, ok := t[key]; ok {

		if m, ok := val.(map[int]int); ok {
			/* act on str */
			return m, nil
		} else {
			return nil, eRROR_NO_GET_WRONG_TYPE
		}
	} else {
		return nil, eRROR_NO_KEY_FOUND
	}
}

func NewTtfInfo() ttfInfo {
	info := make(ttfInfo)
	return info
}
