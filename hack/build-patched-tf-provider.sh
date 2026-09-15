#!/usr/bin/env bash
# Rebuild linux/amd64 patched svalabs/terraform-provider-forgejo for upjet runtime.
#
# Why patch (see docs/fixes.md):
#   - id: number → string (empty external-name before Create)
#   - Read 404 → RemoveResource (Observe → Create)
#   - Delete 404 → success (finalizer cleanup)
#   - TeamMember id = "{team_id}/{user}"
#   - DeployKey key_id / Webhook webhook_id → string; Read skip when unset/0
#
# Usage:
#   ./hack/build-patched-tf-provider.sh
#   TERRAFORM_PROVIDER_VERSION=1.7.0 ./hack/build-patched-tf-provider.sh
#   TERRAFORM_PROVIDER_GIT_REF=abc123def... TERRAFORM_PROVIDER_VERSION=1.7.0-dev \
#     ./hack/build-patched-tf-provider.sh
#
# Env:
#   TERRAFORM_PROVIDER_VERSION  — version string used in output binary name and
#                                 (by default) as git tag v$VERSION (default 1.6.0)
#   TERRAFORM_PROVIDER_GIT_REF  — optional exact git tag/branch/commit to checkout
#                                 (default: v$TERRAFORM_PROVIDER_VERSION)
#   TERRAFORM_PROVIDER_REPO     — git URL (default svalabs upstream)
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
VERSION="${TERRAFORM_PROVIDER_VERSION:-1.6.0}"
GIT_REF="${TERRAFORM_PROVIDER_GIT_REF:-v${VERSION}}"
REPO="${TERRAFORM_PROVIDER_REPO:-https://github.com/svalabs/terraform-provider-forgejo.git}"
OUT_NAME="terraform-provider-forgejo_v${VERSION}"
PATCH_DIR="${ROOT}/hack/svalabs-id-string-patch"
OUT_DIR="${ROOT}/cluster/images/provider-forgejo/patched"
SRC="${TMPDIR:-/tmp}/svalabs-forgejo-build-$$"

cleanup() { rm -rf "$SRC"; }
trap cleanup EXIT

echo "==> clone ${REPO} @ ${GIT_REF} (binary name v${VERSION})"
git clone -c advice.detachedHead=false "$REPO" "$SRC"
git -C "$SRC" fetch --tags --force origin
# Prefer exact ref: tag, branch, or full/short commit SHA
if ! git -C "$SRC" checkout --detach "$GIT_REF" 2>/dev/null; then
  git -C "$SRC" fetch --depth 1 origin "$GIT_REF"
  git -C "$SRC" checkout --detach FETCH_HEAD
fi
echo "==> HEAD $(git -C "$SRC" rev-parse HEAD)"

# Overlays live in this repo with //go:build ignore so golangci-lint / go test
# do not typecheck incomplete svalabs sources against this module. Strip the
# constraint when copying into the svalabs tree so the patched TF provider builds.
overlay() {
  local src="$1" dst="$2"
  # Drop leading build-ignore line(s) and a following blank line if present.
  awk '
    BEGIN { skip = 1 }
    skip && /^\/\/go:build ignore$/ { next }
    skip && /^\/\/ \+build ignore$/ { next }
    skip && /^$/ { skip = 0; next }
    { skip = 0; print }
  ' "$src" > "$dst"
}

echo "==> overlay patched resource files"
mkdir -p "${SRC}/internal/provider"
for f in \
  idstring.go \
  organization_resource.go \
  repository_resource.go \
  team_resource.go \
  team_member_resource.go \
  user_resource.go \
  personal_access_token_resource.go \
  deploy_key_resource.go \
  repository_webhook_resource.go
do
  overlay "${PATCH_DIR}/${f}" "${SRC}/internal/provider/${f}"
done

echo "==> go build linux/amd64"
mkdir -p "$OUT_DIR"
(
  cd "$SRC"
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "${OUT_DIR}/${OUT_NAME}" .
)
chmod +x "${OUT_DIR}/${OUT_NAME}"
echo "==> wrote ${OUT_DIR}/${OUT_NAME}"
file "${OUT_DIR}/${OUT_NAME}"
