package docpdf

// mapOfCharacterToGlyphIndex map of CharacterToGlyphIndex
type mapOfCharacterToGlyphIndex struct {
	keyIndexs map[rune]int //for search index in keys
	Keys      []rune
	Vals      []uint
}

// newMapOfCharacterToGlyphIndex new CharacterToGlyphIndex
func newMapOfCharacterToGlyphIndex() *mapOfCharacterToGlyphIndex {
	var m mapOfCharacterToGlyphIndex
	m.keyIndexs = make(map[rune]int)
	return &m
}

// KeyExists key is exists?
func (m *mapOfCharacterToGlyphIndex) KeyExists(k rune) bool {
	/*for _, key := range m.Keys {
		if k == key {
			return true
		}
	}*/
	if _, ok := m.keyIndexs[k]; ok {
		return true
	}
	return false
}

// Set set key and value to map
func (m *mapOfCharacterToGlyphIndex) Set(k rune, v uint) {
	m.keyIndexs[k] = len(m.Keys)
	m.Keys = append(m.Keys, k)
	m.Vals = append(m.Vals, v)
}

// Index get index by key
func (m *mapOfCharacterToGlyphIndex) Index(k rune) (int, bool) {
	/*for i, key := range m.Keys {
		if k == key {
			return i, true
		}
	}*/
	if index, ok := m.keyIndexs[k]; ok {
		return index, true
	}
	return -1, false
}

// Val get value by Key
func (m *mapOfCharacterToGlyphIndex) Val(k rune) (uint, bool) {
	i, ok := m.Index(k)
	if !ok {
		return 0, false
	}
	return m.Vals[i], true
}

// AllKeys get keys
func (m *mapOfCharacterToGlyphIndex) AllKeys() []rune {
	return m.Keys
}

// AllVals get all values
func (m *mapOfCharacterToGlyphIndex) AllVals() []uint {
	return m.Vals
}
