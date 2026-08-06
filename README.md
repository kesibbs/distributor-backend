# distributor-backend

Go 1.18 distributor API used for an EX288 practice task.

## Practice task

You have been asked to create a reproducible build configuration from the
source code `https://github.com/kesibbs/distributor-backend.git`.

- The application must be built in the project **indigo-build**.
- The created images and the build assets must be named **distributor-backend**.
- The application build must modify the existing S2I scripts of the
  `golang:1.18-ubi9` builder image to set the **SERVER_PORT** environment
  variable with the value **8081**.
- The build must extend from the tag **latest** of the **golang** image
  stream located in the **ocp-images** project (latest tracks `1.18-ubi9`).
- Consider if any permission to reference this image stream tag should be
  given.

The app listens on `SERVER_PORT` (default 8080) and serves `/`, `/health`
and `/api/distributors`. `GET /` reports the active port, so a correct
build answers on 8081.

A full walkthrough is in [SOLUTION.md](SOLUTION.md) — no peeking until
you have attempted it.
