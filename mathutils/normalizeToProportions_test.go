package mathutils

import (
	"testing"
)

func TestNormalizeToProportions(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   []float64
	}{
		{
			name:   "Simple proportions",
			values: []float64{4, 3, 2, 1},
			want:   []float64{0.4, 0.3, 0.2, 0.1},
		},
		{
			name:   "Equal values",
			values: []float64{5, 5, 5, 5},
			want:   []float64{0.25, 0.25, 0.25, 0.25},
		},
		{
			name:   "Single value",
			values: []float64{10},
			want:   []float64{1.0},
		},
		{
			name:   "Zero sum",
			values: []float64{0, 0, 0},
			want:   []float64{0.3333, 0.3333, 0.3334}, // Proporciones iguales con suma 1
		},
		{
			name:   "Empty input",
			values: []float64{},
			want:   []float64{},
		},
		{
			name:   "Negative values",
			values: []float64{-10, -20, -30},
			want:   []float64{0.1667, 0.3333, 0.5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeToProportions(tt.values...)

			// Verificar longitud
			if len(got) != len(tt.want) {
				t.Errorf("NormalizeToProportions() length = %v, want %v", len(got), len(tt.want))
				return
			}

			// Para el caso de suma cero, verificamos que las proporciones sean aproximadamente iguales
			if tt.name == "Zero sum" {
				expectedProportion := 1.0 / float64(len(tt.values))
				for _, v := range got {
					if v < expectedProportion-0.01 || v > expectedProportion+0.01 {
						t.Errorf("NormalizeToProportions() zero sum case should have equal proportions")
						return
					}
				}
				return
			}

			// Verificar valores con tolerancia para floating point
			for i := range got {
				if got[i] < tt.want[i]-0.0001 || got[i] > tt.want[i]+0.0001 {
					t.Errorf("NormalizeToProportions() = %v, want %v", got, tt.want)
					return
				}
			}

			// Verificar que la suma sea cercana a 1.0 (excepto para array vacío)
			if len(got) > 0 {
				var sum float64
				for _, v := range got {
					sum += v
				}
				if sum < 0.999 || sum > 1.001 {
					t.Errorf("NormalizeToProportions() sum = %v, should be close to 1.0", sum)
				}
			}
		})
	}
}
