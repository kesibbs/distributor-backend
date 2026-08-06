# Solution walkthrough (spoilers!)

## 1. Project

```bash
oc login -u developer
oc new-project indigo-build
```

## 2. The S2I customization (already in this repo)

The builder image keeps its original scripts at `/usr/libexec/s2i`
(`STI_SCRIPTS_PATH`). `.s2i/bin/run` overrides the image's run script:
it exports `SERVER_PORT=8081` and then `exec`s the original. Nothing to
do at build time — S2I picks up `.s2i/bin/` automatically.

## 3. Permission for the image stream tag

The `golang` image stream lives in the **ocp-images** project and was
imported with `--reference-policy=local`, so build pods pull it through
the internal registry of that project. The builder service account of
`indigo-build` therefore needs `system:image-puller` there:

```bash
oc policy add-role-to-user system:image-puller \
  system:serviceaccount:indigo-build:builder -n ocp-images
```

Without it the build fails with an image pull authorization error.

## 4. Build configuration

```bash
oc new-app --name distributor-backend \
  ocp-images/golang:latest~https://github.com/kesibbs/distributor-backend.git
oc logs -f bc/distributor-backend      # expect "---> distributor-backend custom assemble/run"
```

`latest` is an alias tag tracking `golang:1.18-ubi9`, so the build
extends the intended builder. All generated assets (bc, is, workload)
are named distributor-backend via `--name`.

## 5. Service and route

`go-toolset` declares no EXPOSE ports, so `oc new-app` cannot create a
Service — create it against the workload with the custom port:

```bash
oc expose deployment/distributor-backend --port=8081   # dc/ if --as-deployment-config
oc expose svc/distributor-backend
```

## 6. Verify

```bash
oc exec deploy/distributor-backend -- env | grep SERVER_PORT   # SERVER_PORT=8081
curl http://$(oc get route distributor-backend -o jsonpath='{.spec.host}')/
# expect "port":"8081" in the JSON
```
