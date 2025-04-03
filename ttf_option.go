package docpdf

// ttfOption  font option
type ttfOption struct {
	UseKerning                bool
	Style                     int               //Regular|Bold|Italic
	OnGlyphNotFound           func(r rune)      //Called when a glyph cannot be found, just for debugging
	OnGlyphNotFoundSubstitute func(r rune) rune //Called when a glyph cannot be found, we can return a new rune to replace it.
}

func defaultTtfFontOption() ttfOption {
	var defa ttfOption
	defa.UseKerning = false
	defa.Style = Regular
	defa.OnGlyphNotFoundSubstitute = defaultOnGlyphNotFoundSubstitute
	return defa
}

func defaultOnGlyphNotFoundSubstitute(r rune) rune {
	return rune('\u0020')
}
