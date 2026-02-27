# 🚀 Roadmap de Escalabilidad: Evolución del IXP Lab hacia una Plataforma SRE Tier-1

El trabajo actual representa una excelente base *Plug & Play* para reemplazar emuladores legacy como Vagrant y VirtualBox por Containerlab. Sin embargo, para llevar el proyecto a un nivel profesional que un Staff SRE exigiría en un **ISP Tier 1**, es necesario transformar el proyecto de un simple *script* de automatización a una plataforma completamente escalable, desacoplada y basada en datos (Data-Driven).

A continuación se detallan las áreas de mejora sugeridas, estructuradas de menor a mayor complejidad arquitectónica:

---

## 🏗️ 1. Arquitectura "Data-Driven" (Desacoplar Código y Datos)

Actualmente, las definiciones de la topología (nodos, direcciones IP, enlaces, ASNs) están estructuradas (hardcodeadas) *dentro* del código de Go (`main.go`).

**El problema:** Si mañana necesitas probar un laboratorio de 30 routers BGP o añadir equipos Arista (`cEOS`) o Nokia (`SRL`), tendrías que tocar el código fuente y recompilar la aplicación de Go.
**La Solución Tier-1:**
- **Externalizar el Modelo Conceptual:** Se debe separar la data en un archivo maestro estructurado (ej. `topology_definition.yaml` o `inventory.json`).
- **Implementación en Go (`Viper`):** Utilizar librerías como `spf13/viper` para que Go ingeste ese archivo externo, popule los *structs* automáticamente en tiempo de ejecución, y renderice el motor Containerlab.
  - *SRE Impact:* Permite que ingenieros de red (que no saben programar) puedan crear complejas topologías modificando un YAML base.

---

## 🧩 2. Autogeneración Dinámica de Archivos TOML (BGP Peers)

Tanto `g2_gobgpd.toml` como `rs_gobgpd.toml` son hoy archivos estáticos. A medida que la topología crezca (ej. inyectar `r5`, `r6`, `r7`), tendrías que escribir a mano las sesiones en el TOML.

**La Solución Tier-1:**
- En Go, utilizar un segundo `text/template` destinado exclusivamente para generar los archivos BGP.
- Go debería iterar por todos los nodos dentro de la L2 del IXP, deducir las IPv4, y autogenerar el array `[[neighbors]]` que requiere el demonio de GoBGP de manera 100% matemática y repetitiva.
  - *SRE Impact:* Garantiza consistencia masiva (Evita "typos" en direcciones o ASNs) y emula el comportamiento de plantillas Jinja/SaltStack que utilizan las *Telcos* reales.

---

## 🛠️ 3. Refactorización de Código Go (Arquitectura MVC/Hexagonal)

El script recae íntegramente en un monobloque dentro de la función `func main()`. La escalabilidad técnica exige aplicar patrones de diseño limpios.

**La Solución Tier-1:**
- Mudar de `main.go` a una estructura profesional modular:
  - `/models` o `/types`: Para albergar la definición de `Node`, `Link`, etc.
  - `/generator`: Logica interna de lectura de `template` y compilación de strings.
  - `/ipam`: Un paquete dedicado a verificar coaliciones de IP.
  - `/cmd`: Donde iría el punto de inicio real de la aplicación, posiblemente utilizando [Cobra](https://github.com/spf13/cobra) para admitir comandos CLI como `ixplab generate --nodes 50`.
  - *SRE Impact:* Testeabilidad. Puedes escribir *Unit Tests* y *Mocking* en la generación de templates o en el manejo de IPs usando `go test`.

---

## 🤖 4. Control del State BGP vía APIs Nativas (gRPC) en vez de Exec CLI

Actualmente el proceso de Go delega en comandos físicos (`ip addr add...` y ejecutar el ejecutable en binario) usando la directiva `Exec` del Containerlab.

**La Solución Tier-1:**
- **Network Programability:** Aprovechar lo que verdaderamente brilla en GoBGP: Su API en lenguajes nativos. En lugar de generar archivos TOML estáticos montados en un contenedor, `go-ixplab` podría actuar como un verdadero Controlador SDN.
- El contenedor arranca en blanco. El código de Go usa cliente **gRPC** (a través de los binarios pre-importados de la librería de GoBGP) y le inyecta directamente vía red las rutas BGP, los prefijos e instruye las modificaciones de *Next-Hops* en el Route Server sin recargar servicios (`Hot Reload`).
  - *SRE Impact:* Este es el estándar hiper-escala adoptado por la industria Cloud (K8s Calico, Cilium, BGP L3).

---

## 🌐 5. Automatización CI/CD "Linterización" Completa

En entornos de Misión Crítica todo debe estar sometido al paradigma *Shift-left testing* (Validar todo en el CI antes del CD).

**La Solución Tier-1:**
- Incorporar **GitHub Actions Workflow** (.github/workflows) que logre:
  1. Ejecutar tests unitarios al código Go (*go fmt*, revisor léxico).
  2. Ejecutar un linter de Containerlab para verificar la sintaxis del YAML resultante.
  3. Comprobar solapamiento de direcciones IPAM Management /10.254.x.x
  4. Levantar la topología usando runners con anidación docker-in-docker, correr el colector `BGP` y tumbarlo nuevamente para garantizar que compila correctamente y no existe *loop de routing*.
  - *SRE Impact:* Imposibilidad de mandar a producción algo que rompa la tabla de control (GitOps Pura).

---

### Resumen para el Portfolio:
Si aplicas estos 5 pasos, pasarás de un **"Excelente Laboratorio Automatizado Plug & Play"** a un **"Framework Open Source de Simulacion de Misión Crítica"** digno de presentaciones en comunidades técnicas Cloud-Native. 

Para comenzar la Fase 2, mi recomendación es atacar el **Punto 1 y 2** simultáneamente, convirtiendo `ixplab` en una verdadera utilería CLI construida sobre Go.
