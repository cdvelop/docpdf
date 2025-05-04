# Especificación: Unificación del Sistema de Manejo de Fuentes en docpdf y chart

## 1. Contexto y Problema

Actualmente, el sistema `docpdf` se enfrenta a un problema de inconsistencia en el manejo de fuentes:

- **docpdf** utiliza [`fontmaker`](fontmaker ) para preprocesar fuentes TrueType y luego renderizarlas mediante su motor interno ([`PdfEngine`](PdfEngine ))
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

- **fontmaker**: Preprocesa fuentes TTF para uso en docpdf
- **PdfEngine**: Motor que renderiza fuentes en PDF
- **chart**: Biblioteca que genera gráficos utilizando freetype
- **[`docFont.go`](docFont.go )**: Define configuraciones de fuentes, pero solo para docpdf. La estructura `FontConfig` contiene la lógica de configuración de fuentes que se desea unificar, pero su ubicación actual crea problemas de dependencias circulares al intentar compartirla con otras partes del sistema.
- **fontbridge**: Paquete actual que intenta compartir configuraciones entre chart y docpdf, pero resulta muy difícil de mantener debido a la duplicidad de código y la lógica fragmentada

### Flujos actuales:

1. **Para texto regular en PDF**:
   - Se cargan fuentes procesadas por fontmaker
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
- [ ] **Etapa 2: Mapeo de dependencias** - Analizar y documentar todas las dependencias actuales de freetype y fontbridge
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

4. **Eliminar dependencia de freetype**:
   - Eliminar importaciones de [`freetype`](freetype )
   - Mantener solo el renderizador SVG
   - Modificar cualquier función que dependa de freetype

5. **Mejorar el renderizador SVG**:
   - Asegurar que puede especificar fuentes por nombre/ruta
   - Implementar todas las características necesarias (tamaños, colores, etc.)

6. **Crear sistema de referencia de fuentes**:
   ```go
   type FontReference struct {
       Name      string
       Weight    string // regular, bold, etc.
       Style     string // normal, italic, etc.
       FontFamily string // Para referenciar en SVG
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
