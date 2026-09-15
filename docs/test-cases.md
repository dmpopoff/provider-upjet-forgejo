<!-- created: 2026-09-15 -->
# provider-forgejo (upjet) — generic test checklist

Create / modify / delete expectations for managed resources generated from
svalabs/forgejo **v1.6.0**. Use against any Forgejo instance + this provider
package (patched TF binary required — see [build.md](./build.md)).

| Legend | Meaning |
|--------|---------|
| **pass** | Synced+Ready; external object matches; delete clears CR **without** stripping finalizers; API confirms |
| **check** | Worth verifying when adding the kind to a cluster |
| **n/a-crd** | TF datasource — no Crossplane CRD; assert via Forgejo HTTP API or `status.atProvider` |

Related: [build.md](./build.md) · [fixes.md](./fixes.md) · examples under `examples/namespaced/`

**Pass rule:** `kubectl delete` must finish with an empty finalizer list; no
manual `kubectl patch` to remove `finalizer.managedresource.crossplane.io`.

---

## 1. Core kinds

| ID | Kind | Create | Modify | Delete | Notes |
|----|------|--------|--------|--------|-------|
| MR-ORG-01 | Organization | check | description/visibility | check | external-name = numeric id |
| MR-REPO-01 | Repository | check | description / private | check | needs owner org/user |
| MR-TEAM-01 | Team | check | description / permission / unitsMap | check | avoid empty unitsMap |
| MR-TM-01 | TeamMember | check | n/a (RequiresReplace) | check | ext-name `teamId/user` |
| MR-USER-01 | User | check | fullName/description | check | require `forProvider.login` |
| MR-PAT-01 | PersonalAccessToken | check | n/a (RequiresReplace) | check | ProviderConfig basic-auth; scopes preserved in TF patch |

## 2. Extended schema

| ID | Kind | Create | Modify | Delete | Notes |
|----|------|--------|--------|--------|-------|
| MR-COL-01 | Collaborator | check | permission | check | ext-name `repoId/user`; no TF `id` |
| MR-DK-01 | DeployKey | check | title/readOnly | check | string `keyId`; Read skips unset/`0` |
| MR-WH-01 | RepositoryWebhook | check | active/url | check | string `webhookId`; same skip |
| MR-GPG-01 | GpgKey | check | — | check | |
| MR-SSH-01 | SshKey | check | — | check | |
| MR-BP-01 | BranchProtection | check | rules | check | |
| MR-OAS/OAV/RAS/RAV | Actions secret/variable | check | value | check | |

`repositoryId` on Collaborator/DeployKey/Webhook is OpenAPI **number** (unquoted YAML).

### DeployKey / Webhook smoke

1. Ready Repository (numeric id).
2. Apply DeployKey (`repositoryId`, `title`, armored `key`, `readOnly`) and/or
   RepositoryWebhook (`type`, `events`, `config.url`).
3. Expect Synced+Ready; `crossplane.io/external-name` numeric string;
   `status.atProvider.keyId` / `webhookId` set.
4. Delete MRs → CR gone; Forgejo API lists no leftover key/hook.

## 3. Lifecycle / negatives

| ID | Case |
|----|------|
| LX-01 | Repo before Org Ready → wait / recover |
| LX-02 | TeamMember before Team Ready (`teamIdRef`) |
| LX-05 | Orphan external + new MR without external-name → Create “already exists” |
| LX-08 | Delete when external already 404 → CR gone, no stuck finalizer |
| LX-10 | Flip Provider package revision → existing MRs stay Healthy |

## 4. ProviderConfig

| ID | Case |
|----|------|
| PC-01 | JSON basic-auth → MRs Ready |
| PC-02 | `api_token` instead of password |
| PC-03 | Wrong password → Synced=False |
| PC-04 | Wrong host → fail, no crash-loop storm |
| PC-06 | PAT create while PC is token-only → expected fail |

## 5. Package

| ID | Case |
|----|------|
| PK-01 | Install from your registry (pull secret if private) |
| PK-04 | Runtime uses **patched** TF binary (not stock svalabs) |
| PK-05 | Stock binary → Create/Delete fail (negative) |

## 6. Datasources (API)

Prefer MR `status.atProvider` or Forgejo HTTP:

`GET /api/v1/orgs/{name}`, `/repos/{owner}/{repo}`, `/teams/{id}`,
`/teams/{id}/members/{user}`, …
