# Upjet-based Crossplane provider for Forgejo

Crossplane provider built with [Upjet](https://github.com/crossplane/upjet) from
[svalabs/terraform-provider-forgejo](https://github.com/svalabs/terraform-provider-forgejo)
**v1.6.0** (all 16 TF resources).

Layout mirrors [provider-upjet-aws](https://github.com/crossplane-contrib/provider-upjet-aws)
(scaffold + `build` submodule).

- Go module: `github.com/crossplane-contrib/provider-upjet-forgejo`
- Package / image name: `provider-forgejo`
- Namespaced MRs: `*.forgejo.forgejo.m.crossplane.io`
- Cluster MRs: `*.forgejo.forgejo.crossplane.io`
- ProviderConfig: `forgejo.m.crossplane.io/v1beta1` (namespaced) /
  `forgejo.crossplane.io/v1beta1` (cluster)
- TF **data sources** are not Crossplane CRDs (upjet does not emit them)

## Credentials

ProviderConfig Secret key `credentials` must be JSON:

```json
{"host":"https://forgejo.example.com","username":"…","password":"…"}
```

Optional: `api_token` instead of username/password (basic-auth required to
*create* PersonalAccessToken MRs).

## Docs

| Doc | Contents |
|-----|----------|
| [docs/build.md](./docs/build.md) | Generate, patch svalabs TF, build, publish xpkg (GHCR / YC CR) |
| [docs/marketplace-upbound.md](./docs/marketplace-upbound.md) | Publish to `xpkg.upbound.io` / marketplace.upbound.io |
| [docs/upgrading.md](./docs/upgrading.md) | Pin upjet / svalabs TF to exact tag or commit |
| [docs/fixes.md](./docs/fixes.md) | Why the TF patch exists (id / 404 / alt identifiers) |
| [docs/test-cases.md](./docs/test-cases.md) | Generic CRUD / lifecycle checklist |

CI: [`.github/workflows/publish-ghcr.yml`](./.github/workflows/publish-ghcr.yml)
builds the patched TF plugin and pushes
`ghcr.io/<owner>/provider-forgejo:<tag>` on `v*` tags or manual dispatch.
[`.github/workflows/ci.yml`](./.github/workflows/ci.yml) includes a free
`vuln-scan` job (`govulncheck` + Trivy fs); all third-party Actions are pinned
by commit digest. Remote: `https://github.com/dmpopoff/provider-upjet-forgejo`.

```bash
make submodules
export PATH="$(go env GOPATH)/bin:$PATH"
./hack/build-patched-tf-provider.sh
make generate
PLATFORM=linux_amd64 VERSION=v0.2.13 make build
# retag / push runtime image + xpkg to your registry — see docs/build.md
```

## Examples

- `examples/namespaced/` — ProviderConfig + one YAML per managed kind
- `examples-generated/` — upjet-generated samples from TF docs
