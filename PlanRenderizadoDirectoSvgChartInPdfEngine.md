# Plan para Renderizado Directo de Gráficos (chart) en docpdf

## 1. Problema Actual

La integración actual de gráficos en `docpdf` presenta varias limitaciones fundamentales, derivadas tanto del flujo de renderizado como de la gestión de dependencias:

**Flujo de Renderizado Ineficiente (SVG/PNG -> Raster):**

1.  **Generación Intermedia:** `docChart.Draw()` utiliza la librería `chart` para generar una representación del gráfico (SVG/PNG) en un buffer.
2.  **Incrustación como Imagen:** Este buffer se trata como una imagen ráster (`docImage`).
3.  **Dibujo en PDF:** La imagen ráster se incrusta en el PDF.

**Desventajas de este flujo:**

*   **Ineficiencia:** Pasos redundantes de codificación/decodificación y procesamiento de imagen.
*   **Pérdida de Calidad Vectorial:** El gráfico final es ráster, perdiendo nitidez al escalar.
*   **Complejidad en `docChart`:** Obliga a manejar formatos de imagen intermedios.

**Gestión de Fuentes Fragmentada y Dependencias Pesadas (Problema de `fontUnitRefactor.md`):**

*   **Inconsistencia:** `docpdf` usa su motor (`pdfEngine`) con fuentes preprocesadas (`fontmaker`), mientras que la librería `chart` (usada por `docChart`) depende de `freetype` para renderizar texto en los gráficos (actualmente como imágenes).
*   **Duplicación:** Gestión de fuentes duplicada e incoherente entre `docpdf` y `chart`.
*   **Dependencia Raster:** `freetype` es una dependencia relativamente pesada, enfocada en renderizado raster.

**Objetivo Estratégico: Uso en Navegador y Tamaño Reducido:**

*   **Meta Futura:** Se desea poder utilizar `docpdf` en entornos de navegador, idealmente compilado con TinyGo para minimizar el tamaño del binario y la latencia.
*   **Conflicto:** Dependencias como `freetype` y la lógica de procesamiento de imágenes ráster son contraproducentes para este objetivo, aumentando significativamente el tamaño del binario y la complejidad.

**Conclusión del Problema:** El enfoque actual no solo es ineficiente y produce gráficos de menor calidad (rasterizados), sino que también introduce dependencias pesadas (`freetype`) y una gestión de fuentes fragmentada que impiden el objetivo estratégico de tener una librería `docpdf` ligera y apta para compilación con TinyGo y uso en navegadores.

**Objetivo de este Plan:** Implementar un renderizado *directo* de los elementos del gráfico (líneas, rectángulos, texto) sobre la página PDF utilizando las primitivas vectoriales de `pdfEngine`. Esto eliminará la necesidad de formatos intermedios (SVG/PNG), la dependencia de `freetype`, y permitirá usar el sistema de fuentes unificado de `pdfEngine`, alineándose con los objetivos de `fontUnitRefactor.md` y la meta de una librería más ligera.

## 2. Solución Propuesta: `pdfRenderer`

La solución consiste en crear un adaptador que traduzca las instrucciones de dibujo de la librería `chart` a las instrucciones de dibujo de `pdfEngine`.

1.  **Implementar `chart.Renderer`:**
    *   Crear un nuevo tipo `pdfRenderer` dentro del paquete `docpdf`.
    *   Este tipo implementará la interfaz `chart.Renderer` (definida en `c:\Users\Cesar\Packages\Internal\docpdf\chart\renderer.go`).
    *   Cada método de la interfaz `chart.Renderer` implementado por `pdfRenderer` (ej. `LineTo(x, y)`, `SetFillColor(c)`, `Circle(r, x, y)`, `Text(body, x, y)`) realizará una llamada al método correspondiente en la instancia de `pdfEngine` (`gp *pdfEngine`).

2.  **Manejo de Coordenadas y Estilos:**
    *   El `pdfRenderer` recibirá la posición de inicio (`x`, `y` en puntos PDF) y las dimensiones deseadas (`width`, `height` en puntos PDF) en la página.
    *   Internamente, `pdfRenderer` deberá realizar las transformaciones necesarias:
        *   **Offset:** Sumar las coordenadas `x`, `y` de inicio (en puntos PDF) a todas las coordenadas recibidas de `chart`. Se asume que `chart` proporcionará coordenadas relativas al origen del gráfico (0,0) y en una escala compatible con puntos PDF.
        *   **Sin Escalado DPI:** El `pdfRenderer` **no** aplicará un factor de escala basado en DPI. La librería `chart` (en su Fase 3 de modificación) deberá realizar sus cálculos de layout y proporcionar sus comandos de dibujo directamente en unidades compatibles con puntos PDF, utilizando las dimensiones (`targetW`, `targetH` en puntos) y DPI proporcionados si es necesario internamente.
    *   Mantendrá el estado del estilo actual (color de trazo, relleno, grosor de línea, fuente, etc., **en puntos PDF**) y lo aplicará a `pdfEngine` antes de ejecutar los comandos de dibujo.

3.  **Integración en `docChart`:**
    *   Se creará un nuevo método `DrawDirect()` en `docChart`.
    *   Este método calculará la posición y dimensiones finales del gráfico en la página PDF.
    *   Instanciará `pdfRenderer` pasándole la instancia de `pdfEngine`, la posición, dimensiones y DPI.
    *   Llamará a un **nuevo método** (que debe ser expuesto por la librería `chart`) en el objeto del gráfico (`chart.BarChart`, `chart.DonutChart`, etc.) que acepte un `chart.Renderer` como argumento (ej. `chartObject.Draw(pdfRenderer)`).
    *   La librería `chart`, al recibir `pdfRenderer`, ejecutará su lógica de dibujo llamando a los métodos de *nuestro* renderer, que a su vez dibujará directamente en el PDF.

## 3. Plan Detallado de Implementación

**Fase 1: Crear el `pdfRenderer` en `docpdf`**

1.  **Crear Archivo:** `c:\Users\Cesar\Packages\Internal\docpdf\pdf_renderer.go`.
2.  **Definir Struct `pdfRenderer`:**
    ```go
    // filepath: c:\Users\Cesar\Packages\Internal\docpdf\pdf_renderer.go
    package docpdf

    import (
        // "image/color" // No es necesario, drawing.Color es suficiente
        // "github.com/cdvelop/docpdf/chart" // Ya no se importa chart directamente aquí
        "github.com/cdvelop/docpdf/chartengine" // Importar la nueva librería interna
        "github.com/cdvelop/docpdf/drawing" // Usaremos drawing.Color
        // "github.com/cdvelop/docpdf/freetype/truetype" // ¡Eliminar esta dependencia!
        // "math" // Para transformaciones
    )

    // pdfRenderer implementa la interfaz chartengine.Renderer
    type pdfRenderer struct {
        pdf         *pdfEngine // Referencia al motor PDF
        offsetX     float64    // Desplazamiento X en puntos PDF
        offsetY     float64    // Desplazamiento Y en puntos PDF
        targetW     float64    // Ancho objetivo en puntos PDF (para referencia de chartengine)
        targetH     float64    // Alto objetivo en puntos PDF (para referencia de chartengine)
        dpi         float64    // DPI (para referencia de chartengine, no para escalado aquí)

        // Estado actual del estilo (en unidades PDF - puntos)
        strokeColor drawing.Color
        fillColor   drawing.Color
        fontColor   drawing.Color
        strokeWidth float64 // En puntos PDF
        fontSize    float64 // En puntos PDF
        fontName    string // Nombre/Identificador de la fuente actual en pdfEngine
        fontStyle   string // Estilo (ej. "B", "I", "BI", "")
        // ... otros estados necesarios
    }

    // Constructor
    func newPdfRenderer(pdf *pdfEngine, x, y, w, h, dpi float64) *pdfRenderer {
        // DPI se pasa para información, no para escalar en este renderer
        if dpi <= 0 {
            dpi = chart.DefaultDPI // O el DPI por defecto de pdfEngine si existe
        }
        return &pdfRenderer{
            pdf:     pdf,
            offsetX: x,
            offsetY: y,
            targetW: w,
            targetH: h,
            dpi:     dpi,
            // Inicializar fontName/fontStyle/fontSize con los valores por defecto de pdfEngine si es posible
            // Los estilos específicos (ChartLabel, ChartAxisLabel) se aplicarán
            // a través de las llamadas SetFont/SetFontSize/SetFontColor desde la librería chart,
            // la cual habrá sido configurada previamente usando FontConfig.
        }
    }

    // --- Implementación de métodos chartengine.Renderer ---
    // (Similar a la implementación anterior, pero usando tipos y métodos
    // definidos en la interfaz chartengine.Renderer)

    // Ejemplo:
    func (r *pdfRenderer) SetStrokeColor(c drawing.Color) {
        r.strokeColor = c
        // r.pdf.SetStrokeColor(c.R, c.G, c.B)
    }
    func (r *pdfRenderer) SetFillColor(c drawing.Color) {
        r.fillColor = c
        // r.pdf.SetFillColor(c.R, c.G, c.B)
    }
    func (r *pdfRenderer) SetLineWidth(width float64) { // Renombrado desde SetStrokeWidth? Verificar interfaz chartengine
        r.strokeWidth = width // Asume puntos PDF
        // r.pdf.SetLineWidth(r.strokeWidth)
    }
    // ... etc para MoveTo, LineTo, Rect, Text, SetFont, SetFontSize ...
    ```
3.  **Completar Implementación:** Implementar *todos* los métodos de la interfaz `chartengine.Renderer`, asegurando la correcta traducción a la API de `pdfEngine`.

**Fase 2: Modificar `docChart` en `docpdf`**

1.  **Añadir Método `DrawDirect`:** Mantener la estructura de `docChart.DrawDirect()` propuesta anteriormente.
2.  **Adaptar `calculatePosition`:** Asegurarse de que funcione como antes.
3.  **Implementar `prepare...Object`:**
    *   Renombrar/Adaptar `prepareBarChartObject` para que cree y configure un objeto `*chartengine.BarChart` en lugar de `*chart.BarChart`.
    *   **Importante:** Configurar los estilos de fuente/texto del `chartengine.BarChart` (para etiquetas, ejes, etc.) usando `c.doc.fontConfig` (mapeando `docpdf.TextStyle` a los campos de estilo de `chartengine`).
4.  **Llamada al Renderizado:** En `DrawDirect`, la llamada será ahora a `renderErr = chartEngineBarChartObject.Draw(pdfRenderer)`. Esto funcionará una vez completada la Fase 3 (nueva).
5.  **Actualizar Posición:** Mantener la lógica como antes.

**Fase 3 (Nueva): Crear Librería Interna `chartengine`**

1.  **Crear Directorio:** Crear un nuevo directorio `chartengine` dentro de `c:\Users\Cesar\Packages\Internal\docpdf\`.
2.  **Definir Interfaz `Renderer`:** Crear `chartengine/renderer.go` y definir la interfaz `Renderer`. Puede ser un subconjunto de `chart.Renderer` para empezar, conteniendo solo los métodos necesarios para un gráfico de barras (ej. `SetStrokeColor`, `SetFillColor`, `SetLineWidth`, `SetFont`, `SetFontSize`, `SetFontColor`, `MoveTo`, `LineTo`, `Rect`, `Text`, `MeasureText`). Asegurarse de que las unidades esperadas (puntos PDF) estén claras en la documentación de la interfaz.
    ```go
    // filepath: c:\Users\Cesar\Packages\Internal\docpdf\chartengine\renderer.go
    package chartengine

    import "github.com/cdvelop/docpdf/drawing"

    // Renderer define la interfaz para dibujar primitivas de gráfico.
    // Se espera que las coordenadas y dimensiones estén en puntos PDF.
    type Renderer interface {
        // Configuración de Estilo (en puntos PDF)
        SetStrokeColor(c drawing.Color)
        SetFillColor(c drawing.Color)
        SetLineWidth(width float64)
        SetFont(name, style string) // Usar nombres/estilos de pdfEngine
        SetFontSize(size float64)   // Tamaño en puntos PDF
        SetFontColor(c drawing.Color)
        // ... otros métodos de estilo (dash array?)

        // Comandos de Dibujo (coordenadas/dimensiones en puntos PDF)
        MoveTo(x, y float64)
        LineTo(x, y float64)
        Rect(x, y, w, h float64) // Dibujar un rectángulo
        Text(text string, x, y float64)
        // ... otros comandos (Circle, Arc, Path?)

        // Medición (en puntos PDF)
        MeasureText(text string) (width float64) // Devolver ancho en puntos

        // Obtener información del contexto
        GetDPI() float64 // DPI informativo
        GetTargetDimensions() (width, height float64) // Dimensiones objetivo en puntos
    }
    ```
3.  **Definir Tipos de Datos:** Crear `chartengine/types.go` (o similar) para definir estructuras como `Value` (para datos de barras), `Style` (para configurar colores, fuentes - usando puntos PDF), etc.
    ```go
    // filepath: c:\Users\Cesar\Packages\Internal\docpdf\chartengine\types.go
    package chartengine
    import "github.com/cdvelop/docpdf/drawing"

    type Style struct {
        FontName    string
        FontStyle   string
        FontSize    float64 // Puntos PDF
        FontColor   drawing.Color
        FillColor   drawing.Color
        StrokeColor drawing.Color
        StrokeWidth float64 // Puntos PDF
        // ... otros estilos
    }

    type Value struct {
        Label string
        Value float64
        Style Style // Estilo para esta barra/etiqueta
    }
    ```
4.  **Implementar `BarChart`:** Crear `chartengine/bar_chart.go`.
    *   Definir la estructura `BarChart` (conteniendo `Title`, `Values`, `BarWidth`, `BarSpacing`, estilos para ejes, etc.).
    *   Implementar el método `Draw(r Renderer) error`. Esta es la lógica central:
        *   Calcular el layout (posiciones X/Y, ancho/alto de barras, ejes) basándose en `Values`, `BarWidth`, `BarSpacing`, y las dimensiones obtenidas de `r.GetTargetDimensions()`. **Todos los cálculos deben hacerse en puntos PDF.**
        *   Iterar sobre las barras:
            *   Configurar el estilo de la barra en el renderer: `r.SetFillColor(bar.Style.FillColor)`, `r.SetStrokeColor(...)`, `r.SetLineWidth(...)`.
            *   Dibujar el rectángulo de la barra: `r.Rect(x, y, w, h)`.
        *   Iterar sobre las etiquetas/valores:
            *   Configurar el estilo de texto: `r.SetFont(...)`, `r.SetFontSize(...)`, `r.SetFontColor(...)`.
            *   Dibujar el texto: `r.Text(bar.Label, x, y)`.
        *   Dibujar ejes (si aplica):
            *   Configurar estilo de línea: `r.SetStrokeColor(...)`, `r.SetLineWidth(...)`.
            *   Dibujar líneas: `r.MoveTo(...)`, `r.LineTo(...)`.
            *   Dibujar etiquetas de ejes: `r.Text(...)`.

## 4. Beneficios Esperados

*   **Rendimiento:** Eliminación de la sobrecarga de generar y procesar formatos intermedios.
*   **Calidad:** Gráficos vectoriales puros en el PDF.
*   **Control Total:** Lógica de gráficos y renderizado dentro del mismo proyecto.
*   **Simplicidad Inicial:** Enfocarse solo en el gráfico de barras para validar la integración.
*   **Sin Dependencia Externa:** No requiere modificar la librería `chart` por ahora.

## 5. Consideraciones y Desafíos

*   **Reimplementación:** Se reimplementará lógica de layout de gráficos de barras (aunque simplificada).
*   **API de `pdfEngine`:** Sigue siendo crucial que `pdfEngine` exponga los métodos primitivos necesarios.
*   **Coordenadas y Unidades:** `chartengine` debe operar consistentemente en puntos PDF.
*   **Fuentes:** El manejo de fuentes depende de la correcta interacción entre `FontConfig`, `docChart`, `chartengine`, `pdfRenderer` y `pdfEngine`.
*   **Expansión Futura:** Si se necesitan más tipos de gráficos, se añadirán a `chartengine` o se podría reconsiderar la integración con la librería externa `chart` una vez validado el modelo.
*   **Completitud:** Implementar los métodos necesarios en `pdfRenderer` para satisfacer la interfaz `chartengine.Renderer` y las necesidades de `chartengine.BarChart`.
