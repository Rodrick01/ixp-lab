# 🚀 Moderna Arquitectura IXP: Containerlab + GoBGP + Observabilidad

![Topología de Laboratorio IXP Containerlab](ixp_containerlab_topology.png)

Este repositorio aloja la solución **Cloud-Native / SRE** para la modernización de la clásica topología de un Punto de Intercambio de Tráfico (IXP). El proyecto reemplaza un entorno heredado, lento y basado en máquinas virtuales completas (VirtualBox/Vagrant), por un orquestador altamente escalable desarrollado en **Go**, capas de red empacadas en contenedores ultraligeros usando **Containerlab**, y un **Stack Integral de Observabilidad Externa**.

Este proyecto está diseñado demostrando las mejores metodologías de un **Site Reliability Engineer (SRE)** aplicado a infraestructuras ISP Tier-1.

---

## 🌟 Características Principales (Enterprise-Grade)

### 1. 🏗️ Infraestructura como Código (IaC) Orquestada con Go
- **Generación Dinámica:** La topología entera (*nodos, interfaces, puertos, volúmenes de configuración BGP y comandos de inyección IP*) es renderizada al vuelo mediante un aplicativo/orquestador escrito 100% en **Go** (`main.go`), usando *structs* fuertemente tipados y la librería nativa `text/template`.
- **Topología Declarativa:** Se emplea un archivo `*.clab.yml` declarativo interpretado por *Containerlab* en milisegundos.

### 2. ⚡ Contenerización Extrema (Bye Hypervisors)
- **GoBGP + FRR Minimalista:** Los routers virtuales de software se compilaron como contenedores base desde Alpine Linux (pocos megas de peso), lo cual arroja tiempos de *boot* casi instantáneos comparados con el provisionado tradicional Unix.
- **Microsegmentación L2:** El conmutador físico de intercambio (IX) se simula a nivel de *bridge* de Linux dentro de Containerlab.

### 3. 🛡️ Plug & Play "Out of The Box" BGP Auto-Provisioning
- Los nodos GoBGP (`g2`) y el Route Server (`rs`) inician con configuraciones pre-calculadas en TOML (`g2_gobgpd.toml` y `rs_gobgpd.toml`).
- Al desplegar la red, el orquestador se encarga de inyectar los direccionamientos IPv4 directamente sobre las veth interfaces a nivel kernel y lanzar los procesos `gobgpd` para establecer las sesiones iBGP y eBGP **sin interacción manual alguna**.

### 4. 📊 Observabilidad SRE Tier-1 Integrada
- Se ha incorporado soporte métrico en formato *Prometheus* de manera nativa mapeando los puertos (ej. `2112`).
- La topología levanta un clúster *out-of-the-box* compuesto por contenedores **Prometheus** (Extracción BGP y recolección de series de tiempo) y **Grafana** (Visualización), exponiendo latencias, advertencias de *route leaks*, fallos de RPKI y la convergencia global de los *peers* del IXP.

---

## 🛠️ Stack Tecnológico

| Componente | Tecnología / Software Implementado |
| :--- | :--- |
| **Generador de Topología** | Go (`text/template`, structs) |
| **Virtualización de Red** | Containerlab (Docker Runtime) |
| **Routing de Alto Rendimiento** | GoBGP (Control Plane), FRR |
| **Router Carrier-Grade** | Juniper vSRX encapsulado vía *vrnetlab* |
| **Telemetry / SRE Metrics** | Prometheus (Time-Series DB) + Grafana |
| **Gestión de Configuración** | Archivos declarativos TOML |

---

## 🚀 Guía Rápida de Despliegue (Quickstart)

### Prerrequisitos
- Entorno Linux (baremetal o VM) con **Docker** activo.
- Compilador de **Go** (1.20+).
- Instalación local de **[Containerlab](https://containerlab.dev/)**.

### Pasos de Ejecución

1. **Generar la Topología mediante Go:**
   Clona este repositorio, navega a la raíz y recompila el archivo YAML maestro:
   ```bash
   go run main.go
   ```
   > Este paso garantiza que las estructuras de datos, las configuraciones *binded* y las directivas *Exec* se unifiquen y escupan el archivo final `ixp-lab.clab.yml`.

2. **Despliegue de la Red a velocidad Containerlab:**
   ```bash
   sudo clab deploy -t ixp-lab.clab.yml
   ```
   > Con este único comando, presenciarás cómo en cuestión de segundos todo tu entorno Tier-1 recobra vida, las interfaces se conectan, los enrutadores hacen BOOT, y BGP comienza la negociación.

3. **Verificar Visibilidad y Tráfico (GoBGP CLI):**
   ```bash
   docker exec -it clab-ixp-lab-g2 gobgp global rib
   # Deberías ver las tablas de rutas aprendidas via el Route Server
   ```

4. **Acceder a Grafana:**
   Abre una pestaña en el navegador y dirígete a `http://localhost:3000` para chequear el estatus saludable de todas las sesiones de interconexión.

---

## 💡 Próximos Pasos (Roadmap del Portfolio)
1. Ingesta de BMP (BGP Monitoring Protocol) hacia Kafka o un servidor temporal para el procesamiento a fondo de grandes volúmenes de tablas mundiales.
2. Inyección de fallas caóticas (Chaos Engineering): Tumbar sesiones o inyectar prefijos inválidos masivos para visualizar el comportamiento de seguridad (RPKI) y la reconvergencia observada en Prometheus.

> **Nota:** Repositorio creado con propósitos educativos de ingeniería SRE demostrando cómo llevar una arquitectura convencional de pruebas de red al moderno paradigma DevOps/GitOps en la era Cloud-Native.
