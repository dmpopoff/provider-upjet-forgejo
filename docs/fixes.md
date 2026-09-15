<!-- created: 2026-09-15 -->
# Why we patch svalabs (upjet / Crossplane)

Engineering notes for running
[svalabs/terraform-provider-forgejo](https://github.com/svalabs/terraform-provider-forgejo)
**v1.6.0** under [upjet](https://github.com/crossplane/upjet). These issues mostly
do **not** show up in plain `terraform apply` with a populated state.

**Full lab chronology** (canary versions, YC publish, Flux, Datastore runbooks):
kept in the consuming infra project’s docs — not in this repository.

Related: [build.md](./build.md) · [test-cases.md](./test-cases.md)

---

## Summary

| Question | Answer |
|----------|--------|
| Bugs in **svalabs** for upjet mode? | **Yes** — empty/`""` numeric `id`, Read/Delete 404 semantics, alt identifiers |
| Bugs in **upjet** itself? | **Not found** — contracts clash with stock Framework schemas |
| Build deltas vs stock upjet template? | Patched TF binary, schema hook, optional `--provenance=false` for some registries |

Patch sources: [`hack/svalabs-id-string-patch/`](../hack/svalabs-id-string-patch/).  
Rebuild: [`hack/build-patched-tf-provider.sh`](../hack/build-patched-tf-provider.sh).  
Schema hook: [`hack/patch-schema-id-string.py`](../hack/patch-schema-id-string.py).

---

## Symptoms and fixes

### A — Create never starts / empty `id`

Upjet puts external-name `""` into TF `id` before Create. Stock svalabs types `id`
as **number** → parse failure.

**Fix:** `id` → string (`FormatInt`) for Org/Repo/Team/User/PAT/… + schema.json
`number`→`string`. External-name: `config.IdentifierFromProvider` (not
`NameAsIdentifier` — that writes the *name* into `id` and breaks Observe).

`forgejo_user`: keep `login` in `forProvider`; do **not** use
`ParameterAsIdentifier("login")` (it still syncs tfstate `id` and can ForceNew).

### B — Observe does not see “gone”

Stock Read on HTTP 404 returns diagnostics error. Upjet treats that as failure,
not absence.

**Fix:** Read 404 → `resp.State.RemoveResource(ctx)`.

### C — Delete stuck on finalizer

Stock Delete on 404 errors → finalizer never clears.

**Fix:** Delete 404 → success.

### D — TeamMember id

Synthetic `id = "{team_id}/{user}"` in the patch; external-name matches.

### E — Team `units` / `units_map`

Forgejo API rejects empty / invalid units. Prefer explicit `unitsMap` (e.g.
`repo.code`) in manifests — not a TF schema bug.

### F — PAT scopes set correlation

After Create, API may return scopes in a different shape → Framework
“planned set element does not correlate”. **Fix:** preserve planned `scopes` in
state on Create/Read (PAT patch).

### G — Collaborator / DeployKey / Webhook without TF `id`

- **Collaborator:** composite external-name `repository_id/user`;
  `NopSetIdentifierArgument`.
- **DeployKey / Webhook:** `key_id` / `webhook_id` as **string**; Read skip when
  unset/`"0"`; `NopSetIdentifierArgument` (Computed-only attrs must not be
  injected into TF *config* — Framework “read-only attribute” error).

---

## Terraform provider issues (vs upjet lifecycle)

| # | Stock behavior | Why it breaks upjet | Patch |
|---|----------------|---------------------|-------|
| T1 | `id` Int64 / number | Empty external-name before Create | string + schema hook |
| T2 | Read 404 → error | No “gone” → no Create | `RemoveResource` |
| T3 | Delete 404 → error | Stuck finalizer | 404 → success |
| T4 | TeamMember weak id | Observe/Delete | synthetic `team_id/user` |
| T5 | empty team units | API 400 | manifest guidance |
| T6 | DK/WH Int64 alt ids | Read ID 0 / read-only inject | string + skip unset; NopSetID |

Upstream svalabs (string ids + 404=gone) would be the long-term fix; until then
vendored patched binary + schema hook.

---

## Not an upjet bug

AsyncOperation / LateInitialize / eventual Ready after modify are expected.
Wrong credentials → Synced=False is expected.
