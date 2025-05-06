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
7. Crear un nuevo paquete llamado `config` donde la estructura `FontConfig` será refactorizada como `TextStyles`, permitiendo su uso sin dependencias circulares. Esto hará que la API se use como `config.TextStyles` en lugar de estar embebida directamente en el paquete docpdf.

## 5. Etapas de Refactorización

Las siguientes etapas se irán completando secuencialmente:

- [ ] **Etapa 1: Migración a SVG** - Eliminar la funcionalidad de renderizado PNG de chart y usar exclusivamente SVG
- [x] **Etapa 2: Mapeo de dependencias** - Analizar y documentar todas las dependencias actuales de freetype y fontbridge
- [ ] **Etapa 3: Nuevo sistema de fuentes** - Diseñar e implementar el sistema unificado de manejo de fuentes
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
   - ✅ Implementar métodos para encadenar la creación de gráficos (ej: `engine.DonutChart().PieChart()`)
   - ✅ Centralizar configuración común (tamaños, DPI, estilos, paleta de colores)

### Fase 2: Modificación de chart

4. **Eliminar dependencia de freetype**: 👈 ESTAMOS AQUÍ
   - ✅ Identificar todos los métodos que usan directamente tipos de freetype
   - ✅ Crear interfaces para abstraer funcionalidades de fuentes. La interfaz clave para esto es `FontProvider`, que se encuentra en el paquete `fontengine` (y se usa como `fontengine.FontProvider`).
   - ✅ Adaptar las implementaciones existentes para usar `fontengine.FontProvider`
   - ⏳ Reemplazar importaciones y usos de [`freetype`](freetype )
   - ⏳ Mantener solo el renderizador SVG
   - ⏳ Implementar un renderizador que cumpla con la interfaz `Renderer` de chart para dibujar directamente en PDF

5. **Implementar abstracción de fontengine.FontProvider**: ✅
   ```go
   // COMPLETADO: Se ha implementado la interfaz fontengine.FontProvider

   // fontengine.FontProvider es una interfaz que abstrae las propiedades necesarias de una fuente
   type fontengine.FontProvider interface {
       // Identificación de la fuente
       Name() string       // Nombre de la fuente
       Family() string     // Familia de la fuente
       
       // Propiedades de estilo
       Weight() string     // Peso: regular, bold, etc.
       Style() string      // Estilo: normal, italic, etc.
       
       // Propiedades para renderizado SVG
       SVGFontID() string  // ID para referenciar en SVG
       
       // Opcionalmente, para sistemas que necesiten la ruta al archivo
       Path() string       // Ruta al archivo de la fuente
   }
   
   // Adaptador transitorio para compatibilidad con código existente
   type TrueTypeFontAdapter struct {
       Font *truetype.Font
       FontName string
       FontFamily string
       FontWeight string
       FontStyle string
       FontPath string
   }
     func NewTrueTypeFontAdapter(font *truetype.Font, name, family, weight, style, path string) fontengine.FontProvider {
       return &TrueTypeFontAdapter{
           Font:       font,
           FontName:   name,
           FontFamily: family,
           FontWeight: weight,
           FontStyle:  style,
           FontPath:   path,
       }
   }
   ```

6. **Crear función GetDefaultFontProvider**: ✅
   ```go
   // GetDefaultFontProvider returns the default font as a fontengine.FontProvider.
   // Esta es la función preferida para el nuevo código que utiliza
   // la abstracción fontengine.FontProvider en lugar de truetype.Font directamente.
   func GetDefaultFontProvider() (fontengine.FontProvider, error) {
       // Primero obtenemos la fuente por el método anterior
       font, err := GetDefaultFont()
       if err != nil {
           return nil, err
       }
       
       // Crear un adaptador para la fuente
       return &TrueTypeFontAdapter{
           Font:       font,
           FontName:   "Roboto-Medium",
           FontFamily: "Roboto",
           FontWeight: "Medium",
           FontStyle:  "normal",
           FontPath:   "", // No necesitamos la ruta para la fuente incorporada
       }, nil
   }
   ```

### Fase 3: Integración con docpdf

7. **Crear paquete config y mover la configuración de fuentes**:
   ```go
   // Mover desde docFont.go al nuevo paquete config
   package config

   // TextStyles representa la configuración de estilos de texto (anteriormente FontConfig)
   type TextStyles struct {
       Family         Font
       Normal         TextStyle
       Header1        TextStyle
       // ...demás campos...
   }
   ```

8. **Extender el nuevo TextStyles para soportar SVG**:
   ```go
   // En el paquete config
   type TextStyles struct {
       // Campos existentes...
       
       // Mapeo de fuentes para SVG
       SVGFontMappings map[string]string
   }
   ```

9. **Crear funciones de inserción SVG en docpdf**:
   ```go
   // Añadir al Document
   func (d *Document) InsertSVG(svg string, x, y float64, width, height float64) error
   ```

10. **Implementar conversión entre estilos de texto**:
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

13. **Actualizar docpdf para integrar SVG**:
    - Implementar inserción y escalado de SVG
    - Asegurar que las referencias a fuentes funcionen correctamente

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

### Próximos Pasos Inmediatos:
1. Implementar un renderizador PDF (PdfRenderer) que implemente la interfaz Renderer. Los métodos de figuras geométricas ya están centralizados en [`geometricFigures.go`](geometricFigures.go ) en el paquete `pdfengine`, pero es necesario refactorizar estos métodos para que coincidan con la interfaz `Renderer` de chart.
2. Crear una función NewPdfRendererProvider que genere un RendererProvider para dibujar directamente en PDF
3. Modificar los métodos Draw de los gráficos para usar el renderizador PDF cuando corresponda
4. Adaptar los estilos en SVG para usar información de fontengine.FontProvider (familia, peso, estilo)
5. Crear el método Document.InsertSVG() para insertar gráficos SVG en el PDF
6. Eliminar completamente el renderizado PNG y la dependencia de freetype

### Avance actual (Mayo 2025):
Hemos completado la primera fase crítica de la refactorización al abstraer la interfaz de fuentes mediante fontengine.FontProvider. Ahora todos los renderizadores utilizan esta interfaz en lugar de depender directamente del tipo *truetype.Font. Las pruebas han sido actualizadas y todo el código compila correctamente con los cambios.
