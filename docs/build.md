<!-- created: 2026-09-15 -->
# Build: provider-upjet-forgejo (linux/amd64)

How to **generate**, **patch**, **build**, and **publish** this Crossplane
Forgejo provider from
[svalabs/terraform-provider-forgejo](https://github.com/svalabs/terraform-provider-forgejo)
via [upjet](https://github.com/crossplane/upjet).

- TF pin: **svalabs/forgejo v1.6.0** (`Makefile`)
- Runtime embeds Terraform **1.5.7** and a **patched** svalabs plugin
  (stock binary is not enough for Crossplane MR lifecycle — see [fixes.md](./fixes.md))
- Target platform for published images: **linux/amd64** (typical Kubernetes)

Related: [fixes.md](./fixes.md) · [test-cases.md](./test-cases.md) ·
[upgrading.md](./upgrading.md) ·
[marketplace-upbound.md](./marketplace-upbound.md) · [README](../README.md)

---

## Quick steps

Prerequisites: Go (≥1.26), Docker (buildx), `crossplane` CLI, network to GitHub +
HashiCorp releases + your OCI registry.

```bash
cd /path/to/provider-upjet-forgejo

make submodules
export PATH="$(go env GOPATH)/bin:$PATH"

# 1) patched svalabs TF binary for linux/amd64 — REQUIRED
./hack/build-patched-tf-provider.sh
# → cluster/images/provider-forgejo/patched/terraform-provider-forgejo_v1.6.0

# 2) schema + docs + codegen
make generate

# 3) controller + runtime image
PLATFORM=linux_amd64 VERSION=v0.2.13 make build
```

Publish (example — any registry; adjust names/tags):

```bash
REGISTRY=ghcr.io/<org>          # or your private registry
VER=0.2.13
IMG=${REGISTRY}/provider-forgejo-runtime:${VER}
PKG=${REGISTRY}/provider-forgejo:v${VER}

# Retag the image produced by make (name varies by BUILD_REGISTRY):
docker images | grep provider-forgejo
docker tag <local-image> "$IMG"

# Some registries reject buildx attestation indexes — disable if needed:
#   docker buildx build ... --provenance=false --sbom=false

docker push "$IMG"

mkdir -p _output
crossplane xpkg build \
  --package-root=package \
  --embed-runtime-image="$IMG" \
  --package-file=_output/provider-forgejo-v${VER}.xpkg

crossplane xpkg push "$PKG" -f _output/provider-forgejo-v${VER}.xpkg
```

Install:

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-forgejo
spec:
  package: ghcr.io/<org>/provider-forgejo:v0.2.13
  # packagePullSecrets: [{name: …}]  # if the registry is private
```

---

## What the pipeline produces

| Artifact | Role |
|----------|------|
| `config/schema.json` | TF schema after id / alt-id string patch |
| `apis/`, `internal/controller/`, `package/crds/` | Generated types + controllers |
| Runtime OCI image | Controller + Terraform + **patched** svalabs plugin |
| `.xpkg` | Crossplane package (CRDs + meta + embedded runtime) |

Stock svalabs v1.6.0 alone is **not** sufficient for Crossplane Observe/Create/Delete
with empty external-name — see [fixes.md](./fixes.md).

---

## Patched TF provider

```bash
./hack/build-patched-tf-provider.sh
```

Overlays `hack/svalabs-id-string-patch/` onto a shallow clone of svalabs v1.6.0 and
writes `cluster/images/provider-forgejo/patched/terraform-provider-forgejo_v1.6.0`.
The runtime Dockerfile downloads the stock zip, then **overwrites** the plugin with
this binary. Missing patched file → image build fails.

After `make generate`, `hack/patch-schema-id-string.py` forces selected schema
attributes (`id`, `key_id`, `webhook_id`, …) from `number` → `string` so CRDs match
the patched plugin.

---

## Publish to GHCR (GitHub Actions)

Workflow [`.github/workflows/publish-ghcr.yml`](../.github/workflows/publish-ghcr.yml)
(all third-party actions pinned by commit SHA):

1. Builds `./hack/build-patched-tf-provider.sh` (required for the runtime image)
2. Runs `make build.all publish` with
   `XPKG_REG_ORGS=ghcr.io/<github.repository_owner>`
3. Pushes `ghcr.io/<owner>/provider-forgejo:<version>`

Triggers: push of tags `v*` (e.g. `v0.2.13`), or **Actions → Publish to GHCR →
Run workflow** with an optional version.

Package visibility: set the GHCR package to public if installs should work
without `packagePullSecrets`. `GITHUB_TOKEN` is enough for push from Actions
(`permissions.packages: write`).

Manual GHCR push (same as the generic example above with
`REGISTRY=ghcr.io/<owner>`):

```bash
echo "$GITHUB_TOKEN" | docker login ghcr.io -u <github-user> --password-stdin
# then docker push + crossplane xpkg push as in Quick steps
```

For **marketplace.upbound.io** (`xpkg.upbound.io`), see
[marketplace-upbound.md](./marketplace-upbound.md) — that path uses an Upbound
**robot** token, not GHCR.

---

## Publish to Yandex Cloud Container Registry (OCI)

[YC Container Registry](https://yandex.cloud/en/docs/container-registry/) is a
standard OCI registry at `cr.yandex/<registry_id>/…`. Any project can use it the
same way as GHCR — this is not Datastore-specific.

### Prerequisites

- [YC CLI](https://yandex.cloud/en/docs/cli/) (`yc`) authenticated to the cloud
- Docker (or compatible client)
- A Container Registry in the folder (create if needed):

```bash
yc container registry create --name provider-forgejo --folder-id <folder_id>
yc container registry list
# note the registry id, e.g. crpxxxxxxxxxxxxxxxx
```

IAM: the principal used for push needs `container-registry.images.pusher`
(or broader editor) on that registry. For cluster pull: `container-registry.images.puller`
plus a dockerconfig JSON pull Secret / Crossplane `ImageConfig`.

### Auth Docker to `cr.yandex`

```bash
yc container registry configure-docker
# writes a credential helper into ~/.docker/config.json for cr.yandex
```

Alternatively: create a service-account key JSON and
`docker login --username json_key --password-stdin cr.yandex < key.json`.

### Push runtime image + Crossplane xpkg

Yandex CR often **rejects** multi-artifact buildx attestation indexes. Prefer
images built with `--provenance=false --sbom=false` (this repo’s image Makefile
already does that on `make build`). If you rebuild the runtime image by hand,
pass the same flags.

```bash
# after: PLATFORM=linux_amd64 VERSION=v0.2.13 make build
REGISTRY=cr.yandex/<registry_id>
VER=0.2.13
IMG=${REGISTRY}/provider-forgejo-runtime:${VER}
PKG=${REGISTRY}/provider-forgejo:v${VER}

docker images | grep provider-forgejo
docker tag <local-image> "$IMG"
docker push "$IMG"

mkdir -p _output
crossplane xpkg build \
  --package-root=package \
  --embed-runtime-image="$IMG" \
  --package-file=_output/provider-forgejo-v${VER}.xpkg

crossplane xpkg push "$PKG" -f _output/provider-forgejo-v${VER}.xpkg
```

Optional manual runtime rebuild (debug):

```bash
cd cluster/images/provider-forgejo
# requires patched/terraform-provider-forgejo_v1.6.0 and bin/linux_amd64/provider
docker buildx build --load \
  --platform linux/amd64 \
  --provenance=false --sbom=false \
  --build-arg TERRAFORM_VERSION=1.5.7 \
  --build-arg TERRAFORM_PROVIDER_SOURCE=svalabs/forgejo \
  --build-arg TERRAFORM_PROVIDER_VERSION=1.6.0 \
  --build-arg TERRAFORM_PROVIDER_DOWNLOAD_NAME=terraform-provider-forgejo \
  --build-arg TERRAFORM_PROVIDER_DOWNLOAD_URL_PREFIX=https://github.com/svalabs/terraform-provider-forgejo/releases/download/v1.6.0 \
  --build-arg TERRAFORM_NATIVE_PROVIDER_BINARY=terraform-provider-forgejo_v1.6.0 \
  -t "${IMG}" \
  .
```

### Install from a private `cr.yandex` registry

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-forgejo
spec:
  package: cr.yandex/<registry_id>/provider-forgejo:v0.2.13
  packagePullPolicy: IfNotPresent
  packagePullSecrets:
    - name: provider-forgejo-pull   # dockerconfigjson in crossplane-system
---
apiVersion: pkg.crossplane.io/v1beta1
kind: ImageConfig
metadata:
  name: provider-forgejo-ycr
spec:
  matchImages:
    - prefix: cr.yandex/
  registry:
    authentication:
      pullSecretRef:
        name: provider-forgejo-pull
```

Create the pull Secret from a YC SA key (or IAM token) the same way you would for
any private OCI registry (`kubernetes.io/dockerconfigjson` with registry host
`cr.yandex`).

---

## Checklist

- [ ] `make submodules`
- [ ] `./hack/build-patched-tf-provider.sh` (linux/amd64 ELF present)
- [ ] `make generate`
- [ ] `PLATFORM=linux_amd64 VERSION=… make build`
- [ ] Runtime image + xpkg pushed to your registry (GHCR / YC CR / Upbound)
- [ ] Cluster: Provider Healthy; smoke ProviderConfig + Organization (see `examples/`)

GHCR via CI: tag `v*` or run **Publish to GHCR**. Upbound Marketplace:
[marketplace-upbound.md](./marketplace-upbound.md).

When bumping upjet or svalabs, follow [upgrading.md](./upgrading.md) instead of
only editing version strings.
