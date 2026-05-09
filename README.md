# Vigtem Sandbox (CLI)

Sandbox ejecutable, sin interfaz web, orientado a análisis de muestras en carpeta con telemetría detallada.

## Características

- **Ejecución por consola**: solicita o recibe por flag la carpeta a analizar.
- **Ingesta aislada**: copia artefactos a un workspace temporal aislado.
- **Análisis estático inicial**:
  - SHA-256 por archivo.
  - Detección de patrones sospechosos (comandos, persistencia, inyección).
- **Simulación de comportamiento**:
  - Simula ejecución de comandos controlados por archivo para trazar flujo.
- **Debugger/Tracer profesional**:
  - Eventos con timestamp en nanosegundos.
  - Traza de archivos, comandos, memoria (direcciones de buffers/hash), y reporte final.
- **Reporte estructurado**:
  - `analysis_report.json` en workspace.

## Arquitectura

- `cmd/sandbox/main.go`: entrada del binario.
- `internal/cli`: interacción y validación de objetivo.
- `internal/engine`: orquestación de ciclo de análisis.
- `internal/analyzer`: firmas y detección de patrones.
- `internal/debugger`: sistema de trazas de alta granularidad.
- `internal/runtime`: probing de direcciones de memoria.
- `internal/config`: configuración por defecto y metadatos host.

## Compilar

```bash
go mod tidy
go build -o vigtem ./cmd/sandbox
```

## Uso

```bash
./vigtem -t /ruta/a/muestras
# o interactivo
./vigtem
```

## Flujo

1. Selecciona carpeta objetivo.
2. Clona los archivos al workspace aislado.
3. Calcula hash y busca patrones sospechosos.
4. Simula comandos de ejecución por archivo para telemetría.
5. Emite trazas detalladas en consola.
6. Escribe `analysis_report.json`.

## Nota importante

Este proyecto es una **base avanzada de sandbox** para investigación y hardening iterativo. Para alcanzar engaño anti-malware de nivel industrial (emulación total de kernel/userland, API hooking profundo, VM introspection, anti-evasión y orquestación multi-OS real), se requiere ampliar componentes de virtualización, instrumentación de syscalls, y análisis dinámico multi-etapa.
