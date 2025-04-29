# Implementación de Formateadores para Gráficos en docpdf

## Objetivo
Mejorar la visualización de gráficos en la biblioteca docpdf añadiendo métodos encadenados para formatear etiquetas y valores numéricos, utilizando las funciones `TruncateName` y `FormatNumber` de la biblioteca `tinystring`.

## Problema
1. Las etiquetas largas en los gráficos de barras se cortan o no se ven bien
2. Los valores numéricos no tienen separación de miles (ej: se ve 1000000 en lugar de 1.000.000)

## Solución implementada

### 1. Creación del paquete `chartutils`
Se creó un nuevo paquete `chartutils` para albergar las funciones de formateo y evitar duplicación de código:

```go
// chartutils/formatters.go
package chartutils

import (
    "github.com/cdvelop/docpdf/chart"
    "github.com/cdvelop/tinystring"
)

// Define tipos para funciones formateadoras
type LabelFormatter func(string) string
type ValueFormatter func(any) string

// Implementa formateadores por defecto
func DefaultLabelFormatter(label string) string { return label }
func DefaultValueFormatter(value any) string { return chart.FloatValueFormatter(value) }

// Formateador basado en TruncateName de tinystring
func TruncateNameLabelFormatter(maxCharsPerWord, maxWidth int) LabelFormatter {
    return func(label string) string {
        return tinystring.Convert(label).TruncateName(maxCharsPerWord, maxWidth).String()
    }
}

// Formateador basado en FormatNumber de tinystring
func FormatNumberValueFormatter(value any) string {
    rawStr := chart.FloatValueFormatter(value)
    return tinystring.Convert(rawStr).FormatNumber().String()
}

// Conversor para integrar con chart.ValueFormatter
func ConvertChartValueFormatter(formatter ValueFormatter) chart.ValueFormatter {
    return func(v any) string {
        return formatter(v)
    }
}
```

### 2. Modificaciones en `docChart.go`

#### 2.1 Nuevos campos añadidos a la estructura `docChart`:
```go
type docChart struct {
    // ... campos existentes ...
    
    labelFormatter chartutils.LabelFormatter // Formateador para etiquetas
    valueFormatter chart.ValueFormatter      // Formateador para valores
    
    // ... resto de campos existentes ...
}
```

#### 2.2 Inicialización de formateadores por defecto:
```go
// En el constructor AddBarChart()
return &docChart{
    // ... otros campos ...
    labelFormatter: chartutils.DefaultLabelFormatter,
    valueFormatter: chart.FloatValueFormatter,
}
```

#### 2.3 Nuevos métodos para configurar formateadores:
```go
// WithLabelFormatter configura un formateador personalizado para etiquetas
func (c *docChart) WithLabelFormatter(formatter chartutils.LabelFormatter) *docChart {
    c.labelFormatter = formatter
    return c
}

// WithValueFormatter configura un formateador personalizado para valores
func (c *docChart) WithValueFormatter(formatter chart.ValueFormatter) *docChart {
    c.valueFormatter = formatter
    return c
}

// WithTruncateNameFormatter usa TruncateName para formatear etiquetas
func (c *docChart) WithTruncateNameFormatter(maxCharsPerWord, maxWidth int) *docChart {
    c.labelFormatter = chartutils.TruncateNameLabelFormatter(maxCharsPerWord, maxWidth)
    return c
}

// WithThousandsSeparator usa FormatNumber para formatear valores con separadores
func (c *docChart) WithThousandsSeparator() *docChart {
    c.valueFormatter = chartutils.ConvertChartValueFormatter(chartutils.FormatNumberValueFormatter)
    return c
}
```

#### 2.4 Aplicar formateadores en el método Draw():
```go
// Aplicar formateadores antes de crear el gráfico
formattedBars := make([]chart.Value, len(c.bars))
for i, bar := range c.bars {
    formattedBars[i] = bar
    if c.labelFormatter != nil {
        formattedBars[i].Label = c.labelFormatter(bar.Label)
    }
}

// Crear el gráfico de barras
barChart := chart.BarChart{
    // ... otros campos ...
    Bars: formattedBars,
}

// Aplicar formateador de valores
if c.valueFormatter != nil {
    barChart.YAxis.ValueFormatter = c.valueFormatter
}
```

### 3. Ejemplo de uso
Se creó un ejemplo en `example/chart_formatters_example.go` que muestra cómo usar estos nuevos métodos:

```go
// Gráfico sin formateo
chartNoFormat := doc.AddBarChart().
    Title("Ventas por Departamento").
    // ... otras configuraciones ...

// Gráfico con formateo
chartWithFormat := doc.AddBarChart().
    Title("Ventas por Departamento").
    // ... otras configuraciones ...
    WithTruncateNameFormatter(3, 15). // Máximo 3 caracteres por palabra, 15 en total
    WithThousandsSeparator()          // Añadir separadores de miles
```

## Pasos pendientes

1. **Pruebas**: Crear pruebas unitarias para los nuevos formateadores y comprobar su funcionamiento.

2. **Optimización**:
   - Evaluar si hay casos en los que el rendimiento pueda mejorarse
   - Considerar una estrategia de caché para evitar formatear repetidamente los mismos valores

3. **Documentación**:
   - Actualizar la documentación principal de la biblioteca docpdf
   - Añadir ejemplos adicionales que muestren diferentes configuraciones de formateo

4. **Posibles mejoras futuras**:
   - Permitir formateo condicional basado en el valor o la posición
   - Añadir más opciones de formateo para valores numéricos (moneda, porcentaje, etc.)
   - Considerar la posibilidad de formateo específico para cada barra individualmente

## Resultados esperados
- Mejor legibilidad de etiquetas largas en gráficos de barras
- Valores numéricos más legibles con separadores de miles
- API consistente con el resto de la biblioteca, usando métodos encadenados

## Conclusión
La implementación de formateadores para gráficos en docpdf mejora significativamente la usabilidad y apariencia de los gráficos de barras, especialmente cuando se trabaja con datos que tienen nombres largos o valores numéricos grandes. La solución se ha diseñado para mantener la compatibilidad con el código existente y seguir el patrón de método encadenado característico de la biblioteca.
