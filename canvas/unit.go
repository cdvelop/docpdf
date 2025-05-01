package canvas

// The units that can be used in the document
const (
	UnitUnset = iota // No units were set, when conversion is called on nothing will happen
	UnitPT           // Points - 1/72 of an inch, traditional unit in PDF documents
	UnitMM           // Millimeters - 1/10 of a centimeter, metric measurement unit
	UnitCM           // Centimeters - 1/100 of a meter, metric measurement unit
	UnitIN           // Inches - Imperial unit equal to 72 points
	UnitPX           // Pixels - screen unit (by default 96 DPI, thus 72/96 = 3/4 point)

	// The math needed to convert units to points
	conversionUnitPT = 1.0
	conversionUnitMM = 72.0 / 25.4
	conversionUnitCM = 72.0 / 2.54
	conversionUnitIN = 72.0
	//We use a dpi of 96 dpi as the default, so we get a conversionUnitPX = 3.0 / 4.0, which comes from 72.0 / 96.0.
	//If you want to change this value, you can change it at Config.ConversionForUnit
	//example: If you use dpi at 300.0
	//pdf.Start(docpdf.Config{PageSize: *docpdf.PageSizeA4, ConversionForUnit: 72.0 / 300.0 })
	conversionUnitPX = 3.0 / 4.0
)
