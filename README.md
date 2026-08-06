# distributor-backend

Go 1.18 distributor API for EX288 custom-S2I drills. Companion to
wholesale-backend, but built with the `golang:1.18-ubi9` builder image
and a **customized S2I process**.

## Endpoints

| Path                | Returns                                             |
|---------------------|-----------------------------------------------------|
| `/`                 | service info incl. `version` and `profile` env vars |
| `/health`           | `{"status":"ok"}`                                   |
| `/api/distributors` | distributor list                                    |

## Custom S2I pieces (`.s2i/`)

The builder image keeps its original scripts at `/usr/libexec/s2i`
(`STI_SCRIPTS_PATH`). This repo overrides them by shipping scripts in
`.s2i/bin/`, each of which sets environment and then `exec`s the original:

- `.s2i/bin/assemble` — exports `CGO_ENABLED=0` and
  `GOFLAGS="-trimpath -buildvcs=false"` for a reproducible build, then
  runs the image's original assemble.
- `.s2i/bin/run` — exports `APP_PROFILE=s2i-custom` at runtime, then
  runs the image's original run script. Visible in `GET /` as `profile`.
- `.s2i/environment` — sets `DIST_VERSION=1.0.0` for both build and
  runtime. Visible in `GET /` as `version`.

## Drill

```bash
# 1. Builder imagestream (OCP 4.22 ships no Go 1.18 — import it):
oc import-image golang:1.18-ubi9 \
  --from=registry.access.redhat.com/ubi9/go-toolset:1.18 --confirm

# 2. Make the latest tag track 1.18-ubi9 (alias, follows the source tag):
oc tag golang:1.18-ubi9 golang:latest --alias=true

# 3. Build + deploy via the latest alias:
oc new-app golang:latest~https://github.com/kesibbs/distributor-backend.git \
  --name distributor-backend

# 4. Expose and verify the custom S2I env vars took effect:
oc expose svc/distributor-backend
curl http://<route-host>/          # expect profile=s2i-custom, version=1.0.0
```
