# qs-sample-app

App de ejemplo para aprender el CI/CD de la plataforma (Tekton en los clusters de `qs-clusters`/`qs-flux`). Tiene **dos
servicios**, cada uno con su `Dockerfile`, que el pipeline construye en paralelo:

| Servicio | Qué hace | Puerto |
|---|---|---|
| `services/api` | HTTP: `/` saluda con el nombre del pod, `/healthz` responde `ok` | 8080 |
| `services/worker` | Proceso en segundo plano: escribe un log cada `INTERVAL` (30s por omisión) | — |

## Contrato de imagen (lo que la plataforma exige a cualquier repo)

1. Un `Dockerfile` por servicio en `services/<nombre>/`; el contexto de build es esa carpeta.
2. Corre como **usuario no root** (aquí `distroless/static:nonroot`, uid 65532).
3. Si sirve HTTP: puerto declarado con `EXPOSE` y un `/healthz` que responde 200.
4. Configuración por **variables de entorno**; logs a **stdout/stderr**.
5. Imágenes base **fijadas por digest** (`@sha256:...`), para que el build sea reproducible.

La plataforma construye, escanea, sube a `ghcr.io/quintana-solutions/qs-sample-app-<servicio>` y firma. Este repo no
publica imágenes por su cuenta.
