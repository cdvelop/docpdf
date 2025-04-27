package docpdf

import (
	"testing"
)

func TestMultiCellWithMaxLines(t *testing.T) {
	// Crear un nuevo documento para pruebas
	doc := NewDocument()

	// Casos de prueba
	testCases := []struct {
		name      string
		text      string
		maxLines  int
		width     float64
		align     position
		expectErr bool
	}{
		{
			name:      "Texto corto una línea",
			text:      "Texto corto",
			maxLines:  2,
			width:     100,
			align:     Left,
			expectErr: false,
		},
		{
			name:      "Texto largo truncado a 1 línea",
			text:      "Este es un texto muy largo que debe truncarse a una sola línea con puntos suspensivos al final",
			maxLines:  1,
			width:     100,
			align:     Left,
			expectErr: false,
		},
		{
			name:      "Texto largo truncado a 2 líneas",
			text:      "Este es un texto muy largo que debe truncarse a dos líneas, con puntos suspensivos al final de la segunda línea para indicar que hay más contenido que no se muestra",
			maxLines:  2,
			width:     100,
			align:     Left,
			expectErr: false,
		},
		{
			name:      "Texto con justificación",
			text:      "Este es un texto largo que debería estar justificado cuando se muestra en múltiples líneas dentro de la celda, para mejorar su apariencia y legibilidad.",
			maxLines:  2,
			width:     100,
			align:     Justify,
			expectErr: false,
		},
		{
			name:      "Texto con valor maxLines inválido",
			text:      "Texto con valor maxLines negativo",
			maxLines:  -1,
			width:     100,
			align:     Left,
			expectErr: false, // No debería dar error, sino usar valor por defecto (1)
		},
	}

	// Probar cada caso
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Configurar opciones de celda
			cellOpt := cellOption{
				Align:  tc.align,
				Border: 0,
			}

			// Ejecutar el método bajo prueba
			err := doc.MultiCellWithOptionAndMaxLines(&Rect{W: tc.width, H: 50}, tc.text, cellOpt, tc.maxLines)

			// Verificar resultado
			if tc.expectErr && err == nil {
				t.Errorf("Se esperaba un error, pero no ocurrió")
			} else if !tc.expectErr && err != nil {
				t.Errorf("No se esperaba error, pero ocurrió: %v", err)
			}

			// Avanzar posición para siguiente prueba
			doc.Br(10)
		})
	}
}
