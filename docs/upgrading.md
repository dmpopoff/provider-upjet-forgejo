<!-- created: 2026-09-15 -->
# Upgrading dependencies

How to pin this provider to an **exact** upstream version of:

1. [crossplane/upjet](https://github.com/crossplane/upjet) (Go module / codegen)
2. [svalabs/terraform-provider-forgejo](https://github.com/svalabs/terraform-provider-forgejo)
   (Terraform provider → schema, docs, runtime plugin)

Also covered: `crossplane/build` git submodule (Makefiles).

Related: [build.md](./build.md) · [fixes.md](./fixes.md)

Current pins (check the tree — do not trust this paragraph alone):

| Dependency | Where pinned | Example |
|------------|--------------|---------|
| upjet | `go.mod` → `github.com/crossplane/upjet/v2` | `v2.4.1-0.20260728103920-4f6e6e10dff2` |
| svalabs forgejo TF | `Makefile` `TERRAFORM_PROVIDER_VERSION` (+ download URL / binary name) | `1.6.0` |
| HashiCorp Terraform CLI | `Makefile` `TERRAFORM_VERSION` | `1.5.7` (must stay **&lt; 1.6**, BSL) |
| build submodule | `.gitmodules` + recorded SHA | `github.com/crossplane/build` |

---

## 1. Update Upjet to a tag / commit

Upjet is a **Go module**. Codegen (`make generate`) and the controller import
`github.com/crossplane/upjet/v2`.

### 1.1 Pick the ref

```bash
# list recent tags
git ls-remote --tags https://github.com/crossplane/upjet.git | tail -20

# or use a commit SHA from https://github.com/crossplane/upjet/commits
UPJET_REF=v2.5.0          # release tag (preferred when available)
# UPJET_REF=4f6e6e10dff2  # exact commit
```

### 1.2 Bump the module

From the provider root:

```bash
# Tag (Go resolves to a semver / pseudo-version automatically):
go get "github.com/crossplane/upjet/v2@${UPJET_REF}"

# Exact commit (same command — go creates a pseudo-version):
# go get github.com/crossplane/upjet/v2@4f6e6e10dff2abcdef...

go mod tidy
```

Confirm in `go.mod`:

```text
github.com/crossplane/upjet/v2 v2.…   # or v2.x.y-0.yyyymmddhhmmss-<sha>
```

Record the resolved line in the PR/commit message (tag **and** pseudo-version /
SHA so others can reproduce).

### 1.3 Align related Crossplane modules (usually required)

Upjet releases are tied to specific `crossplane-runtime` / `crossplane` API
versions. After bumping upjet, if `go mod tidy` or `make generate` fails on
import versions, bump in lockstep, for example:

```bash
go get github.com/crossplane/crossplane-runtime/v2@<compatible-ref>
go get github.com/crossplane/crossplane/apis/v2@<compatible-ref>
go mod tidy
```

Use the versions from the upjet release notes / that upjet commit’s own `go.mod`.

### 1.4 Optional: `build` submodule

Makefiles live in the [`crossplane/build`](https://github.com/crossplane/build)
submodule. Update only when you need newer makelib (not every upjet bump):

```bash
git submodule update --init
cd build
git fetch origin
git checkout <tag-or-sha>    # e.g. a release commit from crossplane/build
cd ..
git add build
```

### 1.5 Regenerate and verify

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
make generate
go build -o /tmp/provider-forgejo ./cmd/provider
PLATFORM=linux_amd64 VERSION=vX.Y.Z make build   # optional full image
```

Expect codegen churn under `apis/`, `internal/controller/`, `package/crds/`.
Review the diff; fix compile breaks in hand-written files (`config/`,
`internal/clients/`, `hack/`).

---

## 2. Update svalabs Terraform provider to a tag / commit

The Forgejo TF provider is pinned in the **Makefile** and consumed in three
places:

| Step | Mechanism |
|------|-----------|
| Schema for upjet codegen | `terraform providers schema` against registry version `TERRAFORM_PROVIDER_VERSION` |
| Resource docs for examples | `pull-docs` clones tag `v$(TERRAFORM_PROVIDER_VERSION)` |
| Runtime plugin | stock zip from GitHub Releases **overwritten** by `./hack/build-patched-tf-provider.sh` |

Patches live in [`hack/svalabs-id-string-patch/`](../hack/svalabs-id-string-patch/)
and **must be rebased** onto the new upstream sources (see [fixes.md](./fixes.md)).

### 2.1 Prefer an exact **release tag** (published to the Terraform Registry)

Example: move from `1.6.0` → `1.7.0`.

1. Edit `Makefile`:

```make
export TERRAFORM_PROVIDER_VERSION ?= 1.7.0
export TERRAFORM_PROVIDER_DOWNLOAD_URL_PREFIX ?= https://github.com/svalabs/terraform-provider-forgejo/releases/download/v$(TERRAFORM_PROVIDER_VERSION)
export TERRAFORM_NATIVE_PROVIDER_BINARY ?= terraform-provider-forgejo_v1.7.0
```

(`TERRAFORM_PROVIDER_SOURCE` / `TERRAFORM_PROVIDER_REPO` usually stay
`svalabs/forgejo` and the GitHub URL.)

2. Rebase patches onto the new tag:

```bash
# fresh upstream tree
rm -rf /tmp/svalabs-forgejo && git clone --branch v1.7.0 --depth 1 \
  https://github.com/svalabs/terraform-provider-forgejo /tmp/svalabs-forgejo

# compare each file under hack/svalabs-id-string-patch/ with
# /tmp/svalabs-forgejo/internal/provider/<same>.go
# re-apply id-string / 404 / PAT scopes / DK+WH changes, then copy back into hack/
```

3. Rebuild patched plugin + regenerate:

```bash
TERRAFORM_PROVIDER_VERSION=1.7.0 ./hack/build-patched-tf-provider.sh
# → cluster/images/provider-forgejo/patched/terraform-provider-forgejo_v1.7.0

rm -rf .work/terraform   # drop cached schema/init if present
make generate
PLATFORM=linux_amd64 VERSION=vX.Y.Z make build
```

4. Run the [test-cases](./test-cases.md) smoke (at least Org/Repo create+delete and
   DeployKey/Webhook if those resources changed).

Document in the commit: **svalabs tag** + **git SHA** of that tag
(`git rev-parse v1.7.0`).

### 2.2 Exact **git commit** (not yet a registry release)

Use when you must track an unreleased SHA.

**A. Runtime patched binary** (supported by the build script):

```bash
export TERRAFORM_PROVIDER_VERSION=1.7.0-dev          # label for output filename only
export TERRAFORM_PROVIDER_GIT_REF=<full-or-short-sha> # exact upstream commit
./hack/build-patched-tf-provider.sh
# binary: cluster/images/provider-forgejo/patched/terraform-provider-forgejo_v1.7.0-dev
```

Set `TERRAFORM_NATIVE_PROVIDER_BINARY` in the Makefile (or environment) to the
**same** `terraform-provider-forgejo_v${TERRAFORM_PROVIDER_VERSION}` name so the
Dockerfile picks it up.

**B. Schema / `make generate`:**

`terraform providers schema` downloads the provider from the **public registry**
using `TERRAFORM_PROVIDER_VERSION`. An unpublished commit is **not** on the
registry.

Options:

1. **Nearest published version** — set `TERRAFORM_PROVIDER_VERSION` to the last
   release, generate schema from that, and manually reconcile CRD diffs against
   the commit’s docs/schema (acceptable only for tiny commits).
2. **Local provider mirror** — build the TF plugin for the host OS from the same
   commit (with or without patches), place it under a
   [filesystem mirror](https://developer.hashicorp.com/terraform/cli/config/config-file#provider-installation)
   matching `registry.terraform.io/svalabs/forgejo/<version>/…`, point
   `TF_CLI_CONFIG_FILE` at a `terraformrc` with `provider_installation.filesystem_mirror`,
   then re-run the schema target / `make generate`.

Until option 2 is automated in this repo, prefer **tag bumps** for routine
upgrades; use commit pins mainly for the **patched runtime binary** while schema
stays on the last release.

**C. Docs pull (`pull-docs`)** still clones `v$(TERRAFORM_PROVIDER_VERSION)`. For
commit-only bumps, either temporarily point `pull-docs` at that SHA by hand or
keep docs from the last tag.

### 2.3 Checklist after a svalabs bump

- [ ] Makefile `TERRAFORM_PROVIDER_VERSION` / `DOWNLOAD_URL_PREFIX` / `NATIVE_PROVIDER_BINARY` agree
- [ ] `hack/svalabs-id-string-patch/` rebased; `./hack/build-patched-tf-provider.sh` succeeds
- [ ] `hack/patch-schema-id-string.py` still covers resources that need string ids / alt ids
- [ ] `make generate` + compile clean
- [ ] Image contains patched binary (not stock-only)
- [ ] CRUD smoke on changed resources

---

## 3. Recording exact versions in git

For reproducible builds, every bump commit should mention:

```text
upjet:     github.com/crossplane/upjet/v2 @ <tag|pseudo-version> (sha <12>)
svalabs:   terraform-provider-forgejo @ vX.Y.Z (sha <12>)  # or GIT_REF=…
terraform: 1.5.7
build:     crossplane/build @ <sha>
```

Optional: add a short `docs/VERSIONS.md` or a table in the release notes — keep
`Makefile` + `go.mod` as the source of truth.
