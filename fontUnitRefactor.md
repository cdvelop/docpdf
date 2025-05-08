# Especificación: Unificación del Sistema de Manejo de Fuentes en docpdf y chart

## 1. Contexto y Problema

Actualmente, el sistema `docpdf` se enfrenta a un problema de inconsistencia en el manejo de fuentes:

- **docpdf** utiliza [`fontengine`](fontengine ) para preprocesar fuentes TrueType y luego renderizarlas mediante su motor interno ([`PdfEngine`](PdfEngine ))
- **chart** (integrado en docpdf) utiliza la biblioteca [`freetype`](freetype ) para renderizar texto en imágenes bitmap
- Esto provoca una gestión de fuentes fragmentada, incoherente y duplicada
- [`docFont.go`](docFont.go ) intenta centralizar la configuración de fuentes, pero no resuelve esta duplicidad

## 2. Objetivo Principal

Crear un sistema unificado de manejo de fuentes que:

1. Elimine la dependencia de [`freetype`](freetype )
2. Utilice exclusivamente gráficos SVG en lugar de PNG para la biblioteca [`chart`](chart )
3. Permita configurar todas las fuentes desde [`docFont.go`](docFont.go ) de manera centralizada
4. Mantenga la consistencia visual en todo el documento PDF generado

## 3. Estado Actual

### Componentes existentes:

- **fontengine**: Preprocesa fuentes TTF para uso en docpdf
- **PdfEngine**: Motor que renderiza fuentes en PDF
- **chart**: Biblioteca que genera gráficos utilizando freetype
- **[`docFont.go`](docFont.go )**: Define configuraciones de fuentes, pero solo para docpdf. La estructura `FontConfig` contiene la lógica de configuración de fuentes que se desea unificar, pero su ubicación actual crea problemas de dependencias circulares al intentar compartirla con otras partes del sistema.
- **fontbridge**: Paquete actual que intenta compartir configuraciones entre chart y docpdf, pero resulta muy difícil de mantener debido a la duplicidad de código y la lógica fragmentada (se necesita eliminar)

### Flujos actuales:

1. **Para texto regular en PDF**:
   - Se cargan fuentes procesadas por fontengine
   - PdfEngine las renderiza directamente en el PDF

2. **Para gráficos**:
   - chart utiliza freetype para renderizar en bitmap
   - Se insertan las imágenes rasterizadas en el PDF

## 4. Solución Propuesta

Aprovechar que [`chart`](chart ) puede generar SVG (vectorial) para unificar el sistema de fuentes:

1. Modificar [`chart`](chart ) para trabajar exclusivamente con SVG
2. Eliminar la dependencia de [`freetype`](freetype ) 
3. Hacer que [`chart`](chart ) referencie las mismas fuentes que utiliza `docpdf`
4. Implementar un sistema para incrustar SVG directamente en documentos PDF
5. Centralizar la gestión de fuentes en una única estructura
6. Eliminar completamente el paquete **fontbridge**, reemplazándolo con el nuevo sistema unificado de manejo de fuentes
7. Crear un nuevo paquete llamado `config` donde la estructura `FontConfig` será refactorizada como `TextStyles` y se implementará una estructura `FontFamily` para gestionar las fuentes, permitiendo su uso sin dependencias circulares. Esto hará que la API se use como `config.TextStyles` y `config.FontFamily` en lugar de estar embebida directamente en el paquete docpdf.

## 5. Etapas de Refactorización

Las siguientes etapas se irán completando secuencialmente:

- [ ] **Etapa 1: Migración a SVG** - Eliminar la funcionalidad de renderizado PNG de chart y usar exclusivamente SVG
- [x] **Etapa 2: Mapeo de dependencias** - Analizar y documentar todas las dependencias actuales de freetype y fontbridge
- [x] **Etapa 3: Nuevo sistema de fuentes** - Diseñar e implementar el sistema unificado de manejo de fuentes ✅
- [ ] **Etapa 4: Integración SVG** - Implementar el sistema de inserción de SVG en docpdf
- [ ] **Etapa 5: Eliminación de fontbridge** - Refactorizar toda funcionalidad que dependa de fontbridge hacia el nuevo sistema
- [ ] **Etapa 6: Pruebas de integración** - Verificar que todos los tipos de gráficos se rendericen correctamente
- [ ] **Etapa 7: Optimización y documentación** - Optimizar rendimiento y crear documentación del nuevo sistema

## 6. Pasos de Implementación Detallados

### Fase 1: Análisis y Preparación

1. **Analizar la estructura de chart**:
   - ✅ Identificar todas las referencias a [`freetype`](freetype )
   - ✅ Localizar el código de renderizado SVG existente
   - ✅ Entender cómo se manejan las fuentes actualmente

2. **Documentar la API de docpdf para manejo de fuentes**:
   - ✅ Mapear cómo se cargan y configuran las fuentes
   - ✅ Identificar cómo se aplican los estilos

3. **Implementar ChartEngine centralizado**:
   - ✅ Crear estructura `ChartEngine` para centralizar la inicialización de fuentes
   - ✅ Implementar inicialización directa sin usar patrones singleton
   - ✅ Implementar métodos para encadenar la creación de gráficos (ej: `engine.Donut().PieChart()`)
   - ✅ Centralizar configuración común (tamaños, DPI, estilos, paleta de colores)

### Fase 2: Modificación de chart

4. **Eliminar dependencia de freetype**:
   - ✅ Identificar todos los métodos que usan directamente tipos de freetype
   - ✅ Crear interfaces para abstraer funcionalidades de fuentes. La interfaz clave para esto es `FontProvider`, que se encuentra en el paquete `fontengine` (y se usa como `fontengine.FontProvider`).
   - ✅ Adaptar las implementaciones existentes para usar `fontengine.FontProvider`
   - ✅ Implementar un renderizador PdfRenderer que cumpla con la interfaz `Renderer` de chart para dibujar directamente en PDF
   - ⏳ Reemplazar importaciones y usos de [`freetype`](freetype )
   - ⏳ Mantener solo el renderizador SVG

5. **Implementar abstracción de fontengine.FontProvider**: ✅
   - Los métodos de `fontengine.FontProvider` para `PdfEngine` se encuentran en `c:\Users\Cesar\Packages\Internal\docpdf\pdfengine\font_provider.go`.

### Fase 3: Integración con docpdf

7. **Crear paquete config y mover la configuración de fuentes**: ✅ COMPLETADO
   - Se ha creado el paquete `config` en `c:\Users\Cesar\Packages\Internal\docpdf\config\`
   - Se ha implementado la estructura `Font` para representar los archivos de fuentes
   - Se ha renombrado `FontConfig` a `TextStyles` y se ha trasladado al nuevo paquete
   - Se han implementado todas las estructuras necesarias para los estilos de texto
   - El nuevo sistema permite una configuración centralizada y coherente de fuentes

8. **Extender el nuevo TextStyles para soportar SVG**:
   ```go
   // En el paquete config
   type TextStyles struct {
       // Campos existentes...
       
       // Mapeo de fuentes para SVG
       SVGFontMappings map[string]string
   }
   ```

9. **Implementar conversión entre estilos de texto**:
    ```go
    // Convertir de TextStyle de docpdf a formato SVG
    func (ts TextStyle) ToSVGStyle() string
    ```

### Fase 4: Sistema Unificado

11. **Implementar manejador centralizado de fuentes**:
    ```go
    // Sistema unificado
    type FontManager struct {
        RegisteredFonts map[string]*FontInfo
        BasePath       string
        
        // Métodos para docpdf y SVG
        GetPDFFontName(style string) string
        GetSVGFontName(style string) string
    }
    ```

12. **Actualizar chart para usar el nuevo sistema**:
    - Modificar chart para generar solo SVG
    - Referenciar fuentes usando el sistema unificado
    - Eliminar cualquier renderizado bitmap

### Fase 5: Pruebas y Refinamiento

14. **Probar con diferentes tipos de gráficos**:
    - Gráficos de barras
    - Gráficos circulares
    - Gráficos de líneas
    - Verificar que las fuentes se rendericen correctamente

15. **Optimizar rendimiento**:
    - Asegurar que los SVG generados no sean excesivamente grandes
    - Implementar cacheo si es necesario

## 7. Consideraciones Técnicas

- **Compatibilidad**: Asegurar que los cambios sean compatibles con el API existente
- **Rendimiento**: Controlar el tamaño de los SVG generados
- **Fallbacks**: Implementar sistemas de respaldo por si una fuente no está disponible
- **Complejidad**: Los gráficos complejos podrían generar SVG muy grandes

## 8. Entregables

1. **Código actualizado** sin dependencia de freetype
2. **Sistema unificado de manejo de fuentes**
3. **Documentación** de uso del nuevo sistema
4. **Ejemplos** de implementación
5. **Pruebas** que demuestren la coherencia del renderizado

## 9. Criterios de Éxito

- Eliminación completa de la dependencia de freetype
- Todos los gráficos se renderizan como SVG
- Las mismas fuentes se utilizan consistentemente en todo el documento
- La configuración de fuentes está centralizada en un solo lugar
- El resultado visual es coherente entre texto normal y gráficos

## 9.1 Política de Unificación de Fuentes

### Directrices Obligatorias:

1. **Eliminación de Fuentes Redundantes**: 
   - ❌ ELIMINAR completamente `roboto.go` y toda referencia a "Roboto" como fuente predeterminada
   - ❌ ELIMINAR todas las funciones `GetDefaultFont()`, `chart.GetDefaultFont()` y similares
   - ❌ PROHIBIDO mantener duplicados de fuentes en diferentes paquetes

2. **Sistema de Fuentes Centralizado**:
   - ✅ TODAS las fuentes deben cargarse UNA SOLA VEZ al inicio (`NewDocument()`)
   - ✅ SOLO usar fuentes ubicadas en el directorio `fonts/` (regular.ttf, bold.ttf, italic.ttf)
   - ✅ `PdfEngine` es el único propietario legítimo de las fuentes cargadas

3. **Comportamiento de Carga de Fuentes**:
   - Cuando se proporciona solo una fuente (como en `docFont_test.go`), esta se usa para todas las variantes
   - Si falta alguna variante (bold o italic), se utiliza la regular como fallback
   - NO existe el concepto de "fuente predeterminada" hardcodeada en el código

4. **Acceso a Fuentes**:
   - Todos los componentes DEBEN obtener referencias a las fuentes a través de `fontengine.FontProvider`
   - `PdfRenderer` DEBE usar las fuentes ya cargadas y registradas en el `PdfEngine`
   - Para el renderizador SVG, usar la información de familia, peso y estilo de `fontengine.FontProvider`

### IMPORTANTE: NUNCA utilizar funciones GetDefaultFont()

Cualquier uso de funciones que devuelvan una "fuente predeterminada" embebida en el código debe ser considerado un error y eliminado. Las fuentes deben provenir exclusivamente de las cargadas explícitamente al iniciar el documento, como se muestra en `docFont_test.go`.

## 10. Estado Actual y Próximos Pasos

### Completado:
- ✅ Definición de la interfaz fontengine.FontProvider que permitirá abstraer las fuentes
- ✅ Implementación de TrueTypeFontAdapter como adaptador transitorio
- ✅ Modificación de SetFont en los renderizadores para usar fontengine.FontProvider
- ✅ Actualización de la estructura Style para soportar fontengine.FontProvider
- ✅ Implementación del método GetFontProvider en Style
- ✅ Creación de la función GetDefaultFontProvider
- ✅ Actualización de todos los archivos de prueba para usar la nueva interfaz
- ✅ Verificación de que todos los códigos compilen sin errores
- ✅ Creación del paquete config con las estructuras TextStyles y FontFamily
- ✅ Migración de la configuración de fuentes de docFont.go al nuevo paquete

### Próximos Pasos Inmediatos:
1. ✅ Implementar un renderizador PDF (PdfRenderer) que implemente la interfaz Renderer.
2. ✅ Crear una función NewPdfRendererProvider que genere un RendererProvider para dibujar directamente en PDF
3. Modificar los métodos Draw de los gráficos para usar el renderizador PDF cuando corresponda
4. Adaptar los estilos en SVG para usar información de fontengine.FontProvider (familia, peso, estilo)
5. Eliminar completamente el renderizado PNG y la dependencia de freetype
6. Implementar la conversión entre TextStyle y estilos SVG

### Avance actual (Mayo 2025):
Hemos completado importantes avances en la refactorización:

1. **Abstracción de fuentes**: Implementamos la interfaz `fontengine.FontProvider` que ahora todos los renderizadores utilizan en lugar de depender directamente del tipo `*truetype.Font`.

2. **PdfRenderer completo**: Hemos implementado exitosamente un renderizador directo a PDF que:
   - Implementa completamente la interfaz `chart.Renderer`
   - Permite trazar gráficos directamente en el PDF sin generar imágenes intermedias
   - Utiliza el mismo motor de renderizado que el resto del documento (PdfEngine)
   - Incorpora manejo adecuado de:
     - Medición precisa de texto usando `MeasureTextWidth`
     - Operaciones de trazo y relleno usando operaciones nativas de PDF
     - Rotación de texto
     - Gestión de paths (caminos) compleja con MoveTo, LineTo, QuadCurveTo, etc.
     - Transformación de coordenadas entre sistemas de referencia

3. **Creación de RendererProvider**: Se ha implementado `NewPdfRendererProvider` que permite a los componentes de chart obtener un renderizador compatible con el motor PDF.

4. **Prueba de integración inicial**: Se ha modificado `docCharts_test.go` para probar exclusivamente el gráfico tipo Donut con renderizado directo a PDF. Esta prueba nos permitirá:
   - Verificar que el `PdfRenderer` funciona correctamente con componentes de gráfico reales
   - Identificar cualquier discrepancia entre la implementación actual y la esperada
   - Establecer un punto de referencia para comparar el rendimiento y la calidad del renderizado
   - Servir como base para próximas pruebas con los demás tipos de gráficos

El siguiente paso clave es modificar los métodos Draw de los componentes gráficos para que utilicen el nuevo renderizador PDF cuando sea apropiado, y eliminar progresivamente la dependencia de freetype.

## 11. Implementación del Renderizador PDF para Chart

Ahora que hemos implementado exitosamente el `PdfRenderer` para trazar gráficos directamente en PDF, detallamos los aspectos clave de esta implementación:

### Prueba de Integración: Gráfico Donut

Hemos configurado una prueba de integración focalizada utilizando exclusivamente el gráfico tipo Donut como caso de prueba. Esta elección estratégica nos permite:

1. **Acotar el alcance**: Concentrarnos en un solo tipo de gráfico simplifica la detección de problemas específicos del renderizado.
  
2. **Comparar implementaciones**: Al usar un gráfico Donut, que combina paths y texto, podemos evaluar todos los aspectos relevantes del renderizador.

3. **Establecer un punto de referencia**: Este gráfico servirá como referencia para validar la correcta implementación y funcionamiento del PdfRenderer.

La prueba se encuentra en el archivo `docCharts_test.go` y genera un único gráfico Donut renderizado directamente a PDF utilizando nuestro nuevo `PdfRenderer`. Esta prueba es el primer paso hacia un testeo más exhaustivo que eventualmente cubrirá todos los tipos de gráficos disponibles.

**Resultados de la prueba inicial (6 de Mayo, 2025):**

- ✅ La prueba se ejecuta exitosamente sin errores
- ✅ El gráfico Donut se renderiza correctamente en el PDF
- ✅ La integración entre `chart.Donut` y `PdfRenderer` funciona adecuadamente
- ✅ Los estilos de texto y colores se aplican como se espera

Este éxito confirma que la arquitectura diseñada para el renderizado directo a PDF es viable y funcional. A continuación, extenderemos las pruebas para incluir los demás tipos de gráficos (barras, líneas, etc.) y validaremos el comportamiento en casos más complejos.

### Estructura del Renderizador

```go
// PdfRenderer implementa la interfaz chart.Renderer y permite dibujar
// directamente en un PDF usando el motor existente de docpdf
type PdfRenderer struct {
    engine       *pdfengine.PdfEngine
    className    string
    dpi          float64
    strokeColor  style.Color
    fillColor    style.Color
    strokeWidth  float64
    dashArray    []float64
    font         fontengine.FontProvider
    fontSize     float64
    fontColor    style.Color
    textRotation float64

    // Coordenadas de la posición actual (para MoveTo, LineTo, etc.)
    currentX float64
    currentY float64

    // Almacena los puntos del path actual para operaciones Close, Stroke, Fill, FillStroke
    pathStartX    float64
    pathStartY    float64
    pathPoints    []Point // Usar un tipo de punto compatible con PdfEngine
    pathClosed    bool
}

// Point representa un punto en un camino (path) para el renderizado
type Point struct {
    X, Y float64
}
```

### Funciones y Métodos Clave

```go
// NewPdfRenderer crea un nuevo renderizador para PDF
func NewPdfRenderer(engine *pdfengine.PdfEngine) *PdfRenderer {
    // La implementación en chart/pdf_renderer.go (revisada) asigna 'engine' directamente al campo 'font'.
    // Esto implica que *pdfengine.PdfEngine debe satisfacer la interfaz fontengine.FontProvider,
    // lo cual es un punto a verificar (relacionado con la TAREA PENDIENTE 2).
    // El comentario original en chart/pdf_renderer.go ("Usar la fuente predeterminada en lugar del engine")
    // sugiere que la intención podría ser obtener una fuente predeterminada DESDE el engine,
    // lo cual difiere de que el engine MISMO sea el FontProvider.
    // Esta documentación ahora refleja la asignación directa como se encuentra en chart/pdf_renderer.go.
    return &PdfRenderer{
        engine:      engine,
        dpi:         96.0,
        strokeColor: style.ColorBlack,
        fillColor:   style.ColorWhite,
        font:        engine, // Asignación directa según chart/pdf_renderer.go.
        strokeWidth: 1.0,    // Valor según chart/pdf_renderer.go.
        fontSize:    10.0,
        fontColor:   style.ColorBlack,
        pathPoints:  []canvas.Point{}, // Asume que Point en chart/pdf_renderer.go es canvas.Point.
    }
}

// NewPdfRendererProvider retorna un RendererProvider que puede generar
// instancias de PdfRenderer
func NewPdfRendererProvider(engine *pdfengine.PdfEngine) RendererProvider {
    return func(width int, height int) (Renderer, error) {
        return NewPdfRenderer(engine), nil
    }
}

// MeasureText utiliza el motor de PDF para calcular las dimensiones del texto
func (r *PdfRenderer) MeasureText(body string) canvas.Box {
    width, err := r.engine.MeasureTextWidth(body)
    if err != nil {
        width = float64(len(body)) * r.fontSize * 0.5
    }
    height := r.fontSize * 1.2
    return canvas.Box{
        Right:  int(width),
        Bottom: int(height),
        Left:   0,
        Top:    0,
    }
}

// Text dibuja texto en el PDF with rotación
func (r *PdfRenderer) Text(body string, x, y float64, options ...TextOption) {
    rotation := 0.0
    for _, option := range options {
        if option.TextRotation != nil {
            rotation = *option.TextRotation
        }
    }
    
    if rotation != 0 {
        r.engine.Rotate(rotation, x, y)
        r.engine.Text(body, x, y, true)
        r.engine.RotateReset()
    } else {
        r.engine.Text(body, x, y, true)
    }
}

// Stroke aplica un trazo al path actual
func (r *PdfRenderer) Stroke() {
    if len(r.pathPoints) > 1 {
        r.engine.Draw()
        r.pathPoints = make([]Point, 0)
        r.pathClosed = false
    }
}

// Fill rellena el path actual
func (r *PdfRenderer) Fill() {
    if len(r.pathPoints) > 1 {
        r.engine.Fill()
        r.pathPoints = make([]Point, 0)
        r.pathClosed = false
    }
}

// FillStroke rellena y traza el path actual
func (r *PdfRenderer) FillStroke() {
    if len(r.pathPoints) > 1 {
        r.engine.FillDraw()
        r.pathPoints = make([]Point, 0)
        r.pathClosed = false
    }
}
```

### Manejo de Paths

Se implementó soporte completo para trazado de paths (caminos) con:

1. **Seguimiento de puntos**:
   - `MoveTo`: Inicia un nuevo path o mueve el punto actual
   - `LineTo`: Añade una línea al path actual
   - `QuadCurveTo`: Añade una curva cuadrática al path

2. **Operaciones de cierre**:
   - `Close`: Cierra el path actual conectando con el punto inicial

3. **Operaciones de renderizado**:
   - `Stroke`: Dibuja el contorno del path
   - `Fill`: Rellena el área del path
   - `FillStroke`: Rellena y dibuja el contorno

### Transformación de coordenadas

El renderizador se encarga de transformar las coordenadas entre el sistema de chart (origen en esquina superior izquierda) y el sistema de coordenadas PDF (origen en esquina inferior izquierda), preservando la coherencia visual del renderizado.

### Integración con PdfEngine

El renderizador aprovecha las funciones existentes en PdfEngine para operaciones como:

- Dibujo de texto con `Text()`
- Medición de texto con `MeasureTextWidth()`
- Rotación con `Rotate()` y `RotateReset()`
- Operaciones de path con `MoveTo()`, `LineTo()`, etc.
- Operaciones de renderizado con `Draw()`, `Fill()` y `FillDraw()`

## 2. `pdfengine` como `fontengine.FontProvider`

**Objetivo:**
Asegurar que `pdfengine` pueda actuar como un proveedor de fuentes (`fontengine.FontProvider`) para ser utilizado en el nuevo sistema unificado de manejo de fuentes.

**Análisis Actual:**
Una revisión del código de `pdfengine.PdfEngine` (en `pdfengine/pdfEngine.go`) indica que la estructura `PdfEngine` **no implementa directamente** la interfaz `fontengine.FontProvider`. `PdfEngine` gestiona las fuentes (ej. a través de métodos como `AddTTFFontData`) y mantiene una referencia a la fuente actualmente seleccionada, que parece ser de tipo `*ttfSubsetObj` (accesible mediante `gp.curr.FontISubset`).

Por lo tanto, la asignación `font: engine` en `chart/pdf_renderer.go` es incorrecta si se espera que `engine` sea un `fontengine.FontProvider`.

**Posible Solución y Dirección:**
La entidad más probable dentro de `pdfengine` que podría implementar (o ser adaptada para implementar) `fontengine.FontProvider` es `*ttfSubsetObj` o una nueva estructura que encapsule la información de una fuente gestionada por `PdfEngine`.

El `PdfEngine` debería entonces exponer un método para obtener un `fontengine.FontProvider` para la fuente actual o una fuente específica. Por ejemplo:
```go
// En pdfengine.PdfEngine
func (gp *PdfEngine) GetCurrentFontProvider() (fontengine.FontProvider, error) {
    if gp.curr.FontISubset == nil {
        return nil, errors.New("no current font selected in PdfEngine")
    }
    // Aquí, se necesitaría que *ttfSubsetObj implemente FontProvider,
    // o se cree un adaptador.
    // Ejemplo conceptual:
    // return NewFontProviderAdapter(gp.curr.FontISubset), nil

    // Asumiendo que *ttfSubsetObj implementará FontProvider directamente o a través de un wrapper:
    // Esta es una suposición que necesita validación y posible implementación.
    // Por ahora, conceptualmente, podría ser algo como:
    fontProviderCandidate := gp.curr.FontISubset // Este objeto necesitaría implementar la interfaz.
    
    // Comprobación (conceptual) si realmente implementa:
    // var _ fontengine.FontProvider = fontProviderCandidate 
    // Si no, se necesita un adaptador o modificar ttfSubsetObj.

    // Placeholder: la siguiente línea requiere que ttfSubsetObj sea o pueda devolver un FontProvider.
    // Esto podría implicar que ttfSubsetObj tenga métodos como Name(), Family(), etc.,
    // o que se construya un adaptador aquí.
    // return fontProviderCandidate, nil // Esto fallará si ttfSubsetObj no es FontProvider
    
    // Se necesitará un trabajo adicional para que ttfSubsetObj (o un adaptador)
    // implemente fontengine.FontProvider. Por ejemplo:
    if provider, ok := gp.curr.FontISubset.(fontengine.FontProvider); ok {
        return provider, nil
    } else {
        // Alternativamente, construir un adaptador aquí si ttfSubsetObj tiene los datos necesarios.
        // return newTtfSubsetFontProviderAdapter(gp.curr.FontISubset), nil
        return nil, errors.New("current font (ttfSubsetObj) does not implement FontProvider and no adapter is in place")
    }
}
```
Esto requeriría que `*ttfSubsetObj` (o un adaptador) sea modificado para implementar todos los métodos de `fontengine.FontProvider` (`Name`, `Family`, `Weight`, `Style`, `SVGFontID`, `Path`). Se debe investigar si `ttfSubsetObj` ya contiene esta información o puede obtenerla fácilmente. La estructura `TtfOption` usada con `AddTTFFontDataWithOption` y los métodos de `ttfSubsetObj` como `GetFamily()` son puntos de partida para esta investigación.

## Pasos Propuestos

1. **Modificar `ttfSubsetObj` para implementar `fontengine.FontProvider`**:
   - Añadir métodos a `ttfSubsetObj` para cumplir con la interfaz `fontengine.FontProvider`.
   - Alternativamente, crear un adaptador que convierta `ttfSubsetObj` a `fontengine.FontProvider`.

2. **Actualizar `PdfEngine` para exponer el método `GetCurrentFontProvider`**:
   - Implementar el método `GetCurrentFontProvider` en `PdfEngine` para devolver la fuente actual como `fontengine.FontProvider`.

3. **Probar la integración**:
   - Verificar que los cambios permiten que `chart` utilice `PdfEngine` como proveedor de fuentes sin problemas de compatibilidad.

4. **Eliminar dependencias obsoletas**:
   - Una vez verificado el funcionamiento, eliminar cualquier referencia o código relacionado con la antigua gestión de fuentes que ya no sea necesario.

5. **Actualizar documentación**:
   - Asegurar que toda la documentación refleje los cambios en la gestión de fuentes y el uso de `PdfEngine` como proveedor de fuentes.

6. **Entrenamiento y adaptación**:
   - Proveer capacitación o guías para los desarrolladores sobre cómo utilizar el nuevo sistema de fuentes y las implicancias de los cambios realizados.

7. **Monitoreo y soporte post-implementación**:
   - Monitorear el sistema en busca de posibles problemas o áreas de mejora después de la implementación de los cambios.
   - Estar preparado para proporcionar soporte y realizar ajustes según sea necesario basado en el feedback de los usuarios y desarrolladores.

8. **Planificación de futuras mejoras**:
   - Con el nuevo sistema de fuentes en funcionamiento, comenzar a planificar posibles mejoras o nuevas funcionalidades que aprovechen la unificación del sistema de manejo de fuentes.

9. **Revisión y lecciones aprendidas**:
   - Realizar una revisión del proceso de implementación para identificar lecciones aprendidas y oportunidades de mejora para futuros proyectos.

10. **Cierre del proyecto**:
    - Una vez completados y validados todos los pasos anteriores, proceder con el cierre formal del proyecto de unificación del sistema de manejo de fuentes, asegurando que toda la documentación esté actualizada y que se haya transferido el conocimiento necesario al equipo de mantenimiento.
