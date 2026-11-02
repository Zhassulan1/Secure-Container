# Secure Container

![CI](https://github.com/Zhassulan1/Secure-Container/actions/workflows/ci-cd.yml/badge.svg)

A small Go health-check service (fasthttp, `GET /check`) behind nginx, packaged as a hardened container. The service itself is trivial; the repo is about the container setup and the CI around it.

## Hardening

Image (`Dockerfile`):
- multi-stage build: the Go binary is built static (`CGO_ENABLED=0`) and stripped (`-ldflags="-w -s"`), the runtime stage is alpine with nginx and that binary
- `apk upgrade` at build time and floating patch tags (`golang:1.27-alpine`, `alpine:3.24`), Dependabot for the rest
- runs as an unprivileged `appuser`
- setuid/setgid bits are removed from every file in the image
- `HEALTHCHECK` requests `/check` through nginx with busybox `wget` (there is no curl in the image)

Runtime (`docker-compose.yml`, `podman_script.sh`):
- all capabilities dropped
- `no-new-privileges`
- seccomp: the engine's default profile. Compose does not override it; `podman_script.sh` passes `/usr/share/containers/seccomp.json` explicitly, which exists only on Podman hosts
- `pids_limit` 1024
- `/tmp` on tmpfs
- restart on failure; Podman also restarts the container when the healthcheck fails

`containers.conf` sets the same defaults (seccomp profile, `no_new_privileges`, pids limit, log size) for all of the user's Podman containers.

## Run

Docker, from source:

```bash
docker compose up -d --build
curl http://localhost:8005/check
```

Published image (CI pushes it to GHCR on every push to `main`):

```bash
docker run -d --name health-check -p 8005:8080 \
  --cap-drop ALL --security-opt no-new-privileges \
  --tmpfs /tmp:mode=1777 --pids-limit 1024 \
  ghcr.io/zhassulan1/health-check:latest
```

Podman (Fedora/RHEL, the script also opens the port in firewalld):

```bash
bash podman_script.sh            # pulls ghcr.io/zhassulan1/health-check:latest
curl http://localhost:8002/check
```

Or build it locally first: `podman build -t ghcr.io/zhassulan1/health-check:latest .`

To apply `containers.conf`, copy it to `~/.config/containers/containers.conf`.

## Check

```bash
docker inspect -f '{{.State.Health.Status}}' health-check   # healthy
docker exec health-check id                                 # uid=1000(appuser)
docker exec health-check grep -E 'Seccomp:|NoNewPrivs|CapEff' /proc/1/status
# Seccomp: 2, NoNewPrivs: 1, CapEff: 0000000000000000
```

With Podman, use `podman` instead of `docker`.

## CI

`.github/workflows/ci-cd.yml` runs on pushes to `main` and `develop` and on pull requests to `main`:

| Job | |
|---|---|
| Build and Test | `go mod verify`, `go vet`, `go test -race` with a coverage report, static build |
| Docker Build | builds the image, starts it with the compose restrictions and waits until the healthcheck passes |
| Security Scan | Trivy on the built image and on the repository (`aquasec/trivy` pinned by digest), SARIF to the Security tab, fails on fixable CRITICAL/HIGH |
| Semgrep SAST | `semgrep --config auto`, SARIF to the Security tab |
| SonarQube | SonarCloud analysis, skipped for dependabot and fork PRs |
| Publish | pushes the scanned image to `ghcr.io/zhassulan1/health-check` (`latest` and commit sha), only on push to `main` |
