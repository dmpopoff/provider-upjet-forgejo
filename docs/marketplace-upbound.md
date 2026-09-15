<!-- created: 2026-09-15 -->
# Publish to marketplace.upbound.io

How to build this provider’s Crossplane package and publish it to the
[Upbound Marketplace](https://marketplace.upbound.io/) registry
(`xpkg.upbound.io`), then optionally list it publicly.

Official docs:

- [Publish packages](https://docs.upbound.io/manuals/marketplace/repositories/publish-packages/)
- [Creating and pushing packages](https://docs.upbound.io/manuals/marketplace/packages/)
- [Authentication (robot tokens)](https://docs.upbound.io/manuals/marketplace/authentication/)

Related local docs: [build.md](./build.md) · [upgrading.md](./upgrading.md)

---

## Prerequisites

- Upbound account and an **organization** (packages live under
  `xpkg.upbound.io/<org>/<repo>`)
- [`up` CLI](https://docs.upbound.io/cli/) and [`crossplane` CLI](https://docs.crossplane.io/latest/cli/)
- Docker (or compatible) client
- Local build of this provider (patched TF + runtime image) — see [build.md](./build.md)

Pushing to `xpkg.upbound.io` requires a **robot token**. Personal API tokens and
`up login` alone return **401** on push. See
[Authenticate to push packages](https://docs.upbound.io/manuals/marketplace/authentication/).

---

## 1. Package metadata (`package/crossplane.yaml`)

Marketplace listing pages render annotations from the package meta. Ensure
`package/crossplane.yaml` has at least:

```yaml
apiVersion: meta.pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-forgejo
  annotations:
    meta.crossplane.io/maintainer: Dmitry Popov <dmpopof@gmail.com>
    meta.crossplane.io/source: https://github.com/dmpopoff/provider-upjet-forgejo
    meta.crossplane.io/license: Apache-2.0
    meta.crossplane.io/description: |
      Crossplane provider for Forgejo (Upjet + svalabs/terraform-provider-forgejo).
    meta.crossplane.io/readme: |
      Manages Forgejo organizations, repositories, teams, deploy keys, webhooks,
      and related resources. Runtime embeds a patched svalabs TF plugin — see docs/.
spec:
  capabilities:
    - SafeStart
```

Rebuild the xpkg after changing annotations.

---

## 2. Create an Upbound repository

Pick an org and repo name (example: org `dmpopoff`, repo `provider-forgejo`).

```bash
up login
up repository create provider-forgejo
up repository list
```

Or create the repository in the Upbound Console → Repositories → Create.

Package address will be:

```text
xpkg.upbound.io/<org>/provider-forgejo:<semver>
```

Marketplace page (after publish / approval):

```text
https://marketplace.upbound.io/providers/<org>/provider-forgejo/
```

---

## 3. Create a robot and grant write access

1. Open `https://accounts.upbound.io/o/<org>/robots` and create a robot.
2. Copy **access ID** and **token** (token is shown once).
3. Assign the robot to a **team** that has **write** on the target repository
   ([Robots](https://docs.upbound.io/manuals/platform/robots/),
   [repository permissions](https://docs.upbound.io/manuals/marketplace/authentication/)).

If `docker-credential-up` is configured for `xpkg.upbound.io`, it can override
`docker login` and cause confusing errors — scope or remove the helper for push
sessions.

```bash
docker login xpkg.upbound.io \
  -u '<robot-access-id>' \
  -p '<robot-token>'
```

---

## 4. Build runtime image + xpkg locally

Same pipeline as [build.md](./build.md). Example with a temporary local tag used
only as the **embedded** runtime (Upbound package embeds this image; the public
install address is the Upbound xpkg tag):

```bash
cd /path/to/provider-upjet-forgejo
make submodules
export PATH="$(go env GOPATH)/bin:$PATH"

./hack/build-patched-tf-provider.sh
make generate

VER=0.2.13
PLATFORM=linux_amd64 VERSION=v${VER} make build

# Retag local runtime image (name from make / BUILD_REGISTRY varies):
docker images | grep provider-forgejo
IMG=localhost/provider-forgejo-runtime:${VER}
docker tag <local-image> "$IMG"

mkdir -p _output
crossplane xpkg build \
  --package-root=package \
  --embed-runtime-image="$IMG" \
  --package-file=_output/provider-forgejo-v${VER}.xpkg
```

Alternatively, after `make build.all` with
`XPKG_REG_ORGS=xpkg.upbound.io/<org>`, use the produced `.xpkg` under `_output/`
if your make/xpkg machinery already embeds the runtime.

The package **version tag must be semver** (e.g. `v0.2.13`). Marketplace only
publishes versions with valid semver tags.

---

## 5. Push to `xpkg.upbound.io`

```bash
ORG=dmpopoff          # your Upbound org
VER=0.2.13

crossplane xpkg push \
  "${ORG}/provider-forgejo:v${VER}" \
  -f _output/provider-forgejo-v${VER}.xpkg
```

Full install reference:

```text
xpkg.upbound.io/dmpopoff/provider-forgejo:v0.2.13
```

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-forgejo
spec:
  package: xpkg.upbound.io/dmpopoff/provider-forgejo:v0.2.13
```

Private repos need `packagePullSecrets` / ImageConfig like any private OCI
registry.

---

## 6. Publish a public Marketplace listing (optional)

Pushing puts the package in the **registry** (pullable). A **public listing** on
[marketplace.upbound.io](https://marketplace.upbound.io/) is separate:

- New repos often require a **one-time Upbound review**.
- Contact Upbound via the `#upbound` channel on Crossplane Slack with:
  - public Git repo URL (`https://github.com/dmpopoff/provider-upjet-forgejo`)
  - Upbound account to list as owner / contact
  - Upbound repository name (`provider-forgejo`)
- With `up` CLI ≥ v0.39.0 you can mark a repo for publish:

```bash
up repository update provider-forgejo --publish
```

Status meanings (publish vs privacy) are documented in
[Publishing public packages](https://docs.upbound.io/manuals/marketplace/repositories/publish-packages/#publishing-public-packages).
Users can still **pull** ACCEPTED / unpublished packages from the registry.

---

## GHCR vs Upbound

| Target | Address | How |
|--------|---------|-----|
| GitHub Container Registry | `ghcr.io/<owner>/provider-forgejo:v…` | Workflow [publish-ghcr.yml](../.github/workflows/publish-ghcr.yml) or manual push in [build.md](./build.md) |
| Upbound Marketplace registry | `xpkg.upbound.io/<org>/provider-forgejo:v…` | This document (`crossplane xpkg push` + robot) |

CI in this repo publishes to **GHCR** by default. Mirroring to Upbound is a
separate `crossplane xpkg push` (or crane copy) using robot credentials — do not
commit robot tokens.

---

## Checklist

- [ ] Annotations in `package/crossplane.yaml`
- [ ] `./hack/build-patched-tf-provider.sh` + `make generate` + build
- [ ] Upbound repo created; robot with write access
- [ ] `docker login xpkg.upbound.io` with robot
- [ ] `crossplane xpkg push <org>/provider-forgejo:v<semver> -f ….xpkg`
- [ ] (Optional) request public listing / `up repository update … --publish`
